// +build linux

package mptcp

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"testing"
)

var (
	// Entries taken from a real MPTCP connections table, used for testing
	testIPv4MPTCPEntry = []byte(" 1: 9C290BF6 4CC0A727  0 E70E8368:0016                         1134B018:BBE8                         01 01 00000000:00000000 15666")
	testIPv6MPTCPEntry = []byte(" 0: F6635734 353F1E98  1 80A80426100000080000000001C07400:1F90 80A80426100000080000000001208902:93A5 01 01 00000000:00000000 39893")
)

// Swap in mock MPTCP lookup function for tests
func init() {
	lookupMPTCPLinux = generateMockLookupMPTCPLinux()
}

// TestLinux_mptcpEnabled verifies that mptcpEnabled properly detects
// multipath TCP functionality on the current Linux system.
func TestLinux_mptcpEnabled(t *testing.T) {
	// Check function result immediately
	enabled, err := mptcpEnabled()
	if err != nil {
		t.Fatal(err)
	}

	// Check if multipath TCP is available by checking for
	// connections table
	_, err = os.Stat(procMPTCP)
	if os.IsNotExist(err) {
		if enabled {
			t.Fatalf("could not find %s, but mptcpEnabled returned true", procMPTCP)
		}

		return
	}

	// Fatal on other errors
	if err != nil {
		t.Fatal(err)
	}

	// Verify multipath TCP is enabled
	if !enabled {
		t.Fatalf("found %s, but mptcpEnabled returned false", procMPTCP)
	}
}

// TestLinux_hostToHex verifies that hostToHex generates the proper hex
// representation of an input IP address string.
func TestLinux_hostToHex(t *testing.T) {
	var tests = []struct {
		host    string
		hexHost string
		err     error
	}{
		// All tests are constants, to ensure test will break if
		// functionality is changed

		// Invalid IP addresses
		{"localhost", "", ErrInvalidIPAddress},
		{"foobar", "", ErrInvalidIPAddress},

		// Valid IPv4 addresses
		{"8.8.4.4", "04040808", nil},
		{"8.8.8.8", "08080808", nil},
		{"10.10.10.10", "0a0a0a0a", nil},
		{"192.168.1.1", "0101a8c0", nil},
		{"255.255.255.0", "00ffffff", nil},

		// Valid IPv6 addresses (not yet implemented)
		{"0000:0000:0000::0000", "", ErrIPv6NotImplemented},
		{"1111:1111:1111::1111", "", ErrIPv6NotImplemented},
		{"2001:4860:4860::8844", "", ErrIPv6NotImplemented},
		{"2001:4860:4860::8888", "", ErrIPv6NotImplemented},
	}

	for i, test := range tests {
		// Convert IP address to hex representation, check results
		hexHost, err := hostToHex(test.host)
		if err != test.err {
			t.Fatalf("[%02d] unexpected err: %v != %v [test: %v]", i, err, test.err, test)
		}

		if hexHost != test.hexHost {
			t.Fatalf("[%02d] unexpected hexHost: %v != %v [test: %v]", i, hexHost, test.hexHost, test)
		}
	}
}

// TestLinux_u16PortToHex verifies that u16PortToHex generates the proper hex
// representation of an input uint16.
func TestLinux_u16PortToHex(t *testing.T) {
	var tests = []struct {
		port    uint16
		hexPort string
	}{
		// All tests are constants, to ensure test will break if
		// functionality is changed
		{0, "0000"},
		{1, "0001"},
		{100, "0064"},
		{1024, "0400"},
		{2123, "084b"},
		{4873, "1309"},
		{8925, "22dd"},
		{65535, "ffff"},
	}

	for i, test := range tests {
		// Convert port to hex representation, check results
		if hexPort := u16PortToHex(test.port); hexPort != test.hexPort {
			t.Fatalf("[%02d] unexpected hexPort: %v != %v [test: %v]", i, hexPort, test.hexPort, test)
		}
	}
}

// TestLinux_mptcpTableReaderLinux verifies that mptcpTableReaderLinux can properly
// parse a Linux MPTCP connections table for entries.
func TestLinux_mptcpTableReaderLinux(t *testing.T) {
	var tests = []struct {
		lines [][]byte
		entry string
		ok    bool
		err   error
	}{
		// Empty file
		{nil, "", false, io.ErrUnexpectedEOF},
		// Invalid header
		{[][]byte{[]byte("foobar")}, "", false, errInvalidMPTCPTable},
		// Header only, no entries
		{[][]byte{mptcpTableHeader}, "", false, nil},
		// Header, bad entry
		{[][]byte{mptcpTableHeader, []byte("foobar")}, "", false, errInvalidMPTCPEntry},
		// Header, not found IPv4 entry
		{[][]byte{mptcpTableHeader, testIPv4MPTCPEntry}, "1134B018:FFFF", false, nil},
		// Header, good IPv4 entry
		{[][]byte{mptcpTableHeader, testIPv4MPTCPEntry}, "1134B018:BBE8", true, nil},
		// Header, good IPv6 entry
		{[][]byte{mptcpTableHeader, testIPv6MPTCPEntry}, "80A80426100000080000000001208902:FFFF", false, nil},
		// Header, good IPv6 entry
		{[][]byte{mptcpTableHeader, testIPv6MPTCPEntry}, "80A80426100000080000000001208902:93A5", true, nil},
	}

	for i, test := range tests {
		// Store input lines in a buffer, appending each with newline
		buf := bytes.NewBuffer(nil)
		for _, l := range test.lines {
			if _, err := buf.Write(append(l, '\n')); err != nil {
				t.Fatal(err)
			}
		}

		// Attempt to check MPTCP table for entry
		ok, err := mptcpTableReaderLinux(buf, test.entry)
		if err != test.err {
			t.Fatalf("[%02d] unexpected err: %v != %v [test: %v]", i, err, test.err, test)
		}

		if ok != test.ok {
			t.Fatalf("[%02d] unexpected ok: %v != %v [test: %v]", i, ok, test.ok, test)
		}
	}
}

// generateMockLookupMPTCPLinux generates a mock Linux MPTCP lookup table, using
// known data.
func generateMockLookupMPTCPLinux() func(string) (bool, error) {
	// Generate lookup table from known hosts and ports
	lookupSet := make(map[string]struct{})
	for host, port := range hostPorts {
		// Convert host to hex
		hexHost, err := hostToHex(host)
		if err != nil {
			if err == ErrIPv6NotImplemented {
				continue
			}

			panic(err)
		}

		// Generate table key with host and port in hex
		key := strings.ToUpper(net.JoinHostPort(hexHost, u16PortToHex(port)))
		lookupSet[key] = struct{}{}
	}

	// Return function which does lookups with mock data
	return func(hexHostPort string) (bool, error) {
		_, ok := lookupSet[hexHostPort]
		return ok, nil
	}
}
