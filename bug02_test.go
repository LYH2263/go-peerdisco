package peerdisco_test

import (
	"testing"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func TestBug02_SnapshotTagsSliceAlias(t *testing.T) {
	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{SelfID: "s", SelfAddr: "127.0.0.1:1", Transport: tr})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.Join(peerdisco.JoinSpec{ID: "b", Addr: "127.0.0.1:3", Tags: []string{"prod"}}); err != nil {
		t.Fatal(err)
	}
	snap := c.Snapshot()
	snap[0].Tags[0] = "x"
	snap2 := c.Snapshot()
	if snap2[0].Tags[0] != "prod" {
		t.Fatalf("snapshot tags aliased: %v", snap2[0].Tags)
	}
}
