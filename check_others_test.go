// +build !linux

package mptcp

import "testing"

// TestOthers_checkMPTCP verifies that checkMPTCP is not implemented on
// platforms other than Linux.
func TestOthers_checkMPTCP(t *testing.T) {
	ok, err := checkMPTCP("localhost", 8080)
	if ok || err != ErrNotImplemented {
		t.Fatalf("checkMPTCP is not implemented, but returned: (%v, %v)", ok, err)
	}
}
