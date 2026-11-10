package v1

import (
	"strings"
	"testing"
)

// Whole-status merge patches are built from possibly stale objects; an unsynced (0) key must be
// omitted so it never overwrites a key that is already persisted.
func TestVpcStatusBytesOmitsUnsyncedTunnelKey(t *testing.T) {
	unsynced, err := (&VpcStatus{Standby: true}).Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(unsynced), "tunnelKey") {
		t.Fatalf("unsynced tunnelKey must be omitted: %s", unsynced)
	}
	synced, err := (&VpcStatus{TunnelKey: 4}).Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(synced), `"tunnelKey":4`) {
		t.Fatalf("synced tunnelKey must be patched: %s", synced)
	}
}
