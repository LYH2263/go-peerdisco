package peerdisco_test

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/audit"
	"example.com/peerdisco/internal/transport"
)

func TestBug09_AuditFileClosed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{
		SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr, AuditPath: path,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Join(peerdisco.JoinSpec{ID: "a", Addr: "127.0.0.1:2", Tags: []string{"t"}}); err != nil {
		t.Fatal(err)
	}
	if err := c.SetMeta("a", map[string]string{"k": "v"}); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	rotated := filepath.Join(dir, "audit.log.1")
	if err := audit.Rotate(path, rotated); err != nil {
		t.Fatalf("rotate after Close: %v", err)
	}
	if _, err := os.Stat(rotated); err != nil {
		t.Fatal(err)
	}
}
