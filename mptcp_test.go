package mptcp

import (
	"errors"
	"net"
	"strconv"
	"testing"
)

const (
	// Known good IPv4 hosts for tests
	ipv4HostOne = "8.8.8.8"
	ipv4HostTwo = "8.8.4.4"

	// Known bad IPv4 hosts for tests
	ipv4BadHostOne = "0.0.0.0"
	ipv4BadHostTwo = "1.1.1.1"

	// Known good IPv6 hosts for tests
	ipv6HostOne = "2001:4860:4860::8888"
	ipv6HostTwo = "2001:4860:4860::8844"

	// Known bad IPv6 hosts for tests
	ipv6BadHostOne = "0000:0000:0000::0000"
	ipv6BadHostTwo = "1111:1111:1111::1111"

	// Known bad IP addresses for tests
	badIPHostOne = "localhost"
	badIPHostTwo = "foobar"
)

// Mapping of known hosts to known ports
var hostPorts = map[string]uint16{
	ipv4HostOne: 2020,
	ipv4HostTwo: 4040,
	ipv6HostOne: 2020,
	ipv6HostTwo: 4040,
}

// TestIsEnabled verifies that IsEnabled returns the same result as its
// underlying implementation.
func TestIsEnabled(t *testing.T) {
	// Check function result immediately
	enabled, err := mptcpEnabled()
	if err != nil {
		t.Fatal(err)
	}

	enabled2, err := IsEnabled()
	if err != nil {
		t.Fatal(err)
	}

	if enabled != enabled2 {
		t.Fatal("mismatch result between IsEnabled and underlying mptcpEnabled")
	}
}

// TestCheck tests the functionality of Check, using a mock lookup
// table, which mocks the true operating system interface.
func TestCheck(t *testing.T) {
	var tests = []struct {
		host string
		port uint16
		ok   bool
		err  error
	}{
		// Invalid IP addresses
		{badIPHostOne, 0, false, ErrInvalidIPAddress},
		{badIPHostTwo, 0, false, ErrInvalidIPAddress},

		// IPv4

		// Invalid hosts, invalid ports
		{ipv4BadHostOne, 8080, false, nil},
		{ipv4BadHostTwo, 6060, false, nil},

		// Valid hosts, invalid ports
		{ipv4HostOne, 1, false, nil},
		{ipv4HostTwo, 10000, false, nil},

		// Invalid hosts, valid ports
		{ipv4BadHostOne, hostPorts[ipv4HostOne], false, nil},
		{ipv4BadHostTwo, hostPorts[ipv4HostTwo], false, nil},

		// Valid hosts, valid ports
		{ipv4HostOne, hostPorts[ipv4HostOne], true, nil},
		{ipv4HostTwo, hostPorts[ipv4HostTwo], true, nil},

		// IPv6 (not yet implemented)

		// Invalid hosts, invalid ports
		{ipv6BadHostOne, 8080, false, ErrIPv6NotImplemented},
		{ipv6BadHostTwo, 6060, false, ErrIPv6NotImplemented},

		// Valid hosts, invalid ports
		{ipv6HostOne, 1, false, ErrIPv6NotImplemented},
		{ipv6HostTwo, 10000, false, ErrIPv6NotImplemented},

		// Invalid hosts, valid ports
		{ipv6BadHostOne, hostPorts[ipv6HostOne], false, ErrIPv6NotImplemented},
		{ipv6BadHostTwo, hostPorts[ipv6HostTwo], false, ErrIPv6NotImplemented},

		// Valid hosts, valid ports
		{ipv6HostOne, hostPorts[ipv6HostOne], false, ErrIPv6NotImplemented},
		{ipv6HostTwo, hostPorts[ipv6HostTwo], false, ErrIPv6NotImplemented},
	}

	for i, test := range tests {
		// Test using expected input, check for expected results
		// Host and port test values are joined here to avoid lots
		// of aggravating formatting on the above test table values
		ok, err := Check(net.JoinHostPort(test.host, strconv.FormatUint(uint64(test.port), 10)))
		if err != test.err {
			t.Fatalf("[%02d] unexpected err: %v != %v [test: %v]", i, err, test.err, test)
		}

		if ok != test.ok {
			t.Fatalf("[%02d] unexpected ok: %v != %v [test: %v]", i, ok, test.ok, test)
		}
	}
}

// TestCheckSplitHostPort verifies that Check properly splits a hostport
// string, and returns the proper results.
func TestCheckSplitHostPort(t *testing.T) {
	var tests = []struct {
		hostport string
		ok       bool
		err      error
	}{
		// Invalid hostport pair
		{"foobar", false, errors.New("missing port in address")},
		{":8080", false, errors.New("invalid IP address")},

		// Invalid port
		{":foo", false, strconv.ErrSyntax},
		{":-1", false, strconv.ErrSyntax},
		{":1000000000", false, strconv.ErrRange},
	}

	for i, test := range tests {
		// Test using expected input, check for expected results
		ok, err := Check(test.hostport)
		if err != nil {
			// Check for network address error
			if aErr, ok := err.(*net.AddrError); ok {
				if aErr.Err != test.err.Error() {
					t.Fatalf("[%02d] unexpected net.AddrErr.Err: %v != %v [test: %v]", i, aErr.Err, test.err, test)
				}

				continue
			}

			// Check for strconv error
			if sErr, ok := err.(*strconv.NumError); ok {
				if sErr.Err != test.err {
					t.Fatalf("[%02d] unexpected strconv.NumError.Err: %v != %v [test: %v]", i, sErr.Err, test.err, test)
				}

				continue
			}

			// Check for string error
			if err.Error() != test.err.Error() {
				t.Fatalf("[%02d] unexpected err: %v != %v [test: %v]", i, err, test.err, test)
			}
		}

		if ok != test.ok {
			t.Fatalf("[%02d] unexpected ok: %v != %v [test: %v]", i, ok, test.ok, test)
		}
	}
}
