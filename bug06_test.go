package peerdisco_test

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug06_PersistFailureNotJoined(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "nope")
	if err := os.WriteFile(bad, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	// PersistPath 指向普通文件作为目录名 → Save 时 MkdirAll/写失败
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{
		SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr,
		PersistPath: filepath.Join(bad, "members.json"),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	err = c.Join(peerdisco.JoinSpec{ID: "p", Addr: "127.0.0.1:5", Tags: []string{"t"}})
	if err == nil {
		t.Fatal("expected persist error")
	}
	if c.AliveCount() != 0 {
		t.Fatalf("member remained after persist fail: %d", c.AliveCount())
	}
}
