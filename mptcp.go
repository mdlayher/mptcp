// Package mptcp provides detection functionality for active, multipath TCP
// connections from a remote client to the current host.  MIT Licensed.
//
// This package is inspired by the original, PHP-based multipath TCP detection
// functions, courtesy of Christoph Paasch and http://multipath-tcp.org/.
package mptcp

import (
	"errors"
	"net"
	"strconv"
)

var (
	// ErrInvalidIPAddress is returned when an invalid IP address is passed to
	// a function.
	ErrInvalidIPAddress = errors.New("invalid IP address")

	// ErrIPv6NotImplemented is returned when an IPv6 address is passed to a
	// function, because IPv6 detection is not yet implemented.
	ErrIPv6NotImplemented = errors.New("IPv6 detection not yet implemented")

	// ErrNotImplemented is returned when MPTCP detection functionality is not
	// implemented for the current operating system.
	ErrNotImplemented = errors.New("not implemented")
)

// IsMPTCP detects if there is an active multipath TCP connection to this machine,
// originating from the input IP address and port pair.  This functionality
// is operating-system dependent, and may not be implemented on all platforms.
//
// If multipath TCP detection is not implemented on the current operating system,
// this function will return ErrNotImplemented.  In addition, other errors may be
// returned on a failed detection.
//
// If multipath TCP detection is implemented on the current operating system,
// this function will return true or false, depending on if the input client host
// and port are using multipath TCP.
func IsMPTCP(host string, port uint16) (bool, error) {
	return checkMPTCP(host, port)
}

// IsMPTCPHostPort behaves like IsMPTCP, but accepts an input host:port string pair,
// as would typically be provided by the RemoteAddr method of a net.Conn.
func IsMPTCPHostPort(hostport string) (bool, error) {
	// Split input hostport pair
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return false, err
	}

	// Convert port into a uint16
	uPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return false, err
	}

	// Check for multipath TCP connectivity
	return IsMPTCP(host, uint16(uPort))
}
