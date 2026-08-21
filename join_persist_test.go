package peerdisco

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/peerdisco/internal/clock"
	"example.com/peerdisco/internal/transport"
)

// newTestCluster constructs a Cluster pointing PersistPath at the given path.
func newTestCluster(t *testing.T, persistPath string) *Cluster {
	t.Helper()
	c, err := New(Options{
		SelfID:      "self",
		SelfAddr:    "127.0.0.1:0",
		MaxMembers:  16,
		PersistPath: persistPath,
		Clock:       clock.NewFake(time.Unix(0, 0)),
		Transport:   &transport.Mock{},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

// persistFailPath returns a PersistPath whose save must fail on every platform:
// the path's directory is a regular file, so os.MkdirAll inside Save cannot
// create it ("not a directory"). This stands in for the real-world read-only /
// unwritable-directory failure the bug report describes.
func persistFailPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("create blocker file: %v", err)
	}
	return filepath.Join(blocker, "state.json")
}

// TestJoinRollsBackOnPersistFailure reproduces the bug: when persist fails the
// member must NOT remain in the in-memory table. After a failed Join the member
// must be absent (no half-success) so memory and disk stay consistent.
func TestJoinRollsBackOnPersistFailure(t *testing.T) {
	c := newTestCluster(t, persistFailPath(t))
	defer c.Close()

	err := c.Join(JoinSpec{ID: "nodeA", Addr: "10.0.0.1:7946"})
	if err == nil {
		t.Fatalf("Join succeeded despite unwritable persist path; want persist error")
	}
	if !errors.Is(err, ErrPersist) {
		t.Fatalf("Join error = %v, want error wrapping ErrPersist", err)
	}

	// The bug: member appeared in memory even though persist failed.
	if _, ok := c.Member("nodeA"); ok {
		t.Errorf("Member(nodeA) present after failed Join; want rolled back (memory/disk consistency)")
	}
	if got := c.AliveCount(); got != 0 {
		t.Errorf("AliveCount = %d after failed Join, want 0", got)
	}
	if n := len(c.Members()); n != 0 {
		t.Errorf("Members() = %d after failed Join, want 0", n)
	}

	// A subsequent Join to the same ID must not be blocked by a phantom entry.
	if _, ok := c.table.Get("nodeA"); ok {
		t.Fatalf("phantom nodeA still in member table; rollback incomplete")
	}
}

// Control: with a writable persist path, Join succeeds and persists to disk,
// and a fresh cluster reloading from disk sees the same membership.
func TestJoinPersistsWhenPathWritable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	c := newTestCluster(t, path)
	defer c.Close()

	if err := c.Join(JoinSpec{ID: "nodeA", Addr: "10.0.0.1:7946"}); err != nil {
		t.Fatalf("Join: %v", err)
	}
	if got := c.AliveCount(); got != 1 {
		t.Fatalf("AliveCount = %d, want 1", got)
	}

	// Disk should reflect the member.
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("persist file not written: %v", err)
	}

	// A fresh cluster loading from disk should see nodeA (memory/disk consistent).
	c2 := newTestCluster(t, path)
	defer c2.Close()
	if err := c2.LoadPersist(); err != nil {
		t.Fatalf("LoadPersist: %v", err)
	}
	if got := c2.AliveCount(); got != 1 {
		t.Errorf("after reload AliveCount = %d, want 1", got)
	}
	if _, ok := c2.Member("nodeA"); !ok {
		t.Errorf("reloaded cluster missing nodeA")
	}
}
