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

// IsEnabled returns whether or the current host supports multipath TCP.
// If multipath TCP is enabled on this host, this function will return true.
// If it is not enabled on this host, or an error occurs, this function will
// return false.
//
// It is recommended to check the result of IsEnabled before attempting to check
// for active multipath TCP connections using Check.
func IsEnabled() (bool, error) {
	return mptcpEnabled()
}

// Check detects if there is an active multipath TCP connection to this machine,
// originating from the input host:port string, such as one returned from the
// RemoteAddr method of a net.Conn.
//
// This functionality is operating-system dependent, and may not be implemented
// on all platforms. It is recommended to check the result of IsEnabled before
// attempting to check for active multipath TCP connections using Check.
//
// If multipath TCP detection is not implemented for the current operating system,
// this function will return ErrNotImplemented.  In addition, other errors may be
// returned on a failed detection.
//
// If multipath TCP detection is implemented on the current operating system,
// this function will return true or false, depending on if a connection with
// the input host:port string is active and is using multipath TCP.
func Check(hostport string) (bool, error) {
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
	return checkMPTCP(host, uint16(uPort))
}
