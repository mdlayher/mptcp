// +build linux

package mptcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

// procMPTCP is the location of the Linux-specific file which contains
// the active MPTCP connections table.
const procMPTCP = "/proc/net/mptcp"

// checkMPTCP checks if an input host string and uint16 port are present
// in this Linux machine's MPTCP active connections.
var checkMPTCP = func(host string, port uint16) (bool, error) {
	// Get hex representation of host
	hexHost, err := hostToHex(host)
	if err != nil {
		return false, err
	}

	// Combine hex host and port, convert to uppercase
	hexHostPort := strings.ToUpper(net.JoinHostPort(hexHost, u16PortToHex(port)))

	// Use lookup function to check for results
	return lookupMPTCPLinux(hexHostPort)
}

// hostToHex converts an input host IP address into its equivalent hex form,
// for use with MPTCP connection lookup.
func hostToHex(host string) (string, error) {
	// Parse IP address from host
	ip := net.ParseIP(host)

	// If result is not nil, we assume this is IPv4
	if ip4 := ip.To4(); ip4 != nil && len(ip4) == net.IPv4len {
		// For IPv4, grab the IPv4 hex representation of the address
		return fmt.Sprintf("%02x%02x%02x%02x", ip4[3], ip4[2], ip4[1], ip4[0]), nil
	}

	// Check for IPv6 address
	if ip6 := ip.To16(); ip6 != nil && len(ip6) == net.IPv6len {
		// TODO(mdlayher): attempt to check for IPv6 address
		return "", ErrIPv6NotImplemented
	}

	// IP address is not valid
	return "", ErrInvalidIPAddress
}

// u16PortToHex converts an input uint16 port into its equivalent hex form,
// for use with MPTCP connection lookup.
func u16PortToHex(port uint16) string {
	// Store uint16 in buffer using little endian byte order
	portBuf := [2]byte{}
	binary.LittleEndian.PutUint16(portBuf[:], port)

	// Retrieve hex representation of uint16 port
	return fmt.Sprintf("%02x%02x", portBuf[1], portBuf[0])
}

// lookupMPTCPLinux uses the Linux /proc filesystem to attempt to detect
// active MPTCP connections matching the input hex host:port pair.
//
// This implementation is swappable for testing with a mock data source.
var lookupMPTCPLinux = func(hexHostPort string) (bool, error) {
	// Read in the entire MPTCP connections table
	// TODO(mdlayher): consider using a more efficient method (bufio.NewScanner)
	mptcpBuf, err := ioutil.ReadFile(procMPTCP)
	if err != nil {
		return false, err
	}

	// Check if the combined host:port pair are present in the MPTCP
	// connections table
	return bytes.Contains(mptcpBuf, []byte(hexHostPort)), nil
}
