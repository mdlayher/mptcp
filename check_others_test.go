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

// TestOthers_mptcpEnabled verifies that mptcpEnabled always returns
// false unless a platform explicitly supports it.
func TestOthers_mptcpEnabled(t *testing.T) {
	ok, err := mptcpEnabled()
	if ok || err != nil {
		t.Fatalf("mptcpEnabled should return (false, nil), but returned: (%v, %v)", ok, err)
	}
}
