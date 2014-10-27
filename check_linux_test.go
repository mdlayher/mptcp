// +build linux

package mptcp

import (
	"log"
	"net"
	"strings"
	"testing"
)

// Swap in mock MPTCP lookup function for tests
func init() {
	lookupMPTCPLinux = generateMockLookupMPTCPLinux()
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

		log.Printf("mock: %s:%d -> %s", host, port, key)
	}

	// Return function which does lookups with mock data
	return func(hexHostPort string) (bool, error) {
		_, ok := lookupSet[hexHostPort]
		return ok, nil
	}
}
