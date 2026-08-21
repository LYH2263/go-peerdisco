package peerdisco

import (
	"testing"
	"time"

	"example.com/peerdisco/internal/clock"
	"example.com/peerdisco/internal/transport"
)

func newTestCluster(t *testing.T) *Cluster {
	t.Helper()
	c, err := New(Options{
		SelfID:      "self",
		SelfAddr:    "127.0.0.1:1",
		Transport:   &transport.Mock{},
		Clock:       clock.NewFake(time.Unix(0, 0)),
		PersistPath: "", // 不落盘，保持单测独立
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

// TestMembersTagsNotAliased 复现并锁定“前端脱敏写穿集群成员表”的回归：
// Members() 返回的 MemberView.Tags 必须是独立副本，调用方就地把
// Tags[0] 改成 "mask" 不得污染内部 table，否则 gossip 带出脏标签、
// 对端策略匹配全乱。
func TestMembersTagsNotAliased(t *testing.T) {
	c := newTestCluster(t)
	defer c.Close()

	if err := c.Join(JoinSpec{
		ID:   "n1",
		Addr: "127.0.0.1:2",
		Tags: []string{"prod", "cache"},
	}); err != nil {
		t.Fatalf("Join: %v", err)
	}

	views := c.Members()
	if len(views) != 1 {
		t.Fatalf("Members len = %d, want 1", len(views))
	}
	if len(views[0].Tags) != 2 || views[0].Tags[0] != "prod" {
		t.Fatalf("Tags = %v, want [prod cache]", views[0].Tags)
	}

	// 前端脱敏：就地改快照里的首标签。
	views[0].Tags[0] = "mask"
	if got := views[0].Tags[0]; got != "mask" {
		t.Fatalf("快照就地写未生效：%q", got)
	}

	// 再次拉 Members：真实标签必须保持 prod。
	again := c.Members()
	if again[0].Tags[0] != "prod" {
		t.Fatalf("快照写穿了成员表：Tags[0] = %q, want %q", again[0].Tags[0], "prod")
	}

	// Member(id) 同样不能被污染。
	m, ok := c.Member("n1")
	if !ok {
		t.Fatal("Member(n1) 未找到")
	}
	if m.Tags[0] != "prod" {
		t.Fatalf("Member.Tags[0] = %q, want %q", m.Tags[0], "prod")
	}
}

// TestSnapshotTagsNotAliased 对齐 Snapshot 路径的同一回归。
func TestSnapshotTagsNotAliased(t *testing.T) {
	c := newTestCluster(t)
	defer c.Close()

	if err := c.Join(JoinSpec{
		ID:   "n1",
		Addr: "127.0.0.1:2",
		Tags: []string{"prod"},
	}); err != nil {
		t.Fatalf("Join: %v", err)
	}

	snap := c.Snapshot()
	if snap[0].Tags[0] != "prod" {
		t.Fatalf("Snapshot Tags[0] = %q, want prod", snap[0].Tags[0])
	}
	snap[0].Tags[0] = "mask"

	again := c.Snapshot()
	if again[0].Tags[0] != "prod" {
		t.Fatalf("Snapshot 写穿了成员表：Tags[0] = %q, want %q", again[0].Tags[0], "prod")
	}
}
