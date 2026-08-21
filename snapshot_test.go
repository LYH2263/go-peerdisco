package peerdisco

import (
	"sort"
	"testing"

	"example.com/peerdisco/internal/transport"
)

// newTestCluster 构造一个最小可用的 Cluster，用于快照隔离测试。
func newTestCluster(t *testing.T) *Cluster {
	t.Helper()
	c, err := New(Options{
		SelfID:   "self",
		SelfAddr: "127.0.0.1:1",
		Transport: &transport.Mock{},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

// TestSnapshotTagsAliasing 复现巡检脚本污染问题：
// 第一次 Snapshot 拿到 Tags 后本地排序并改写一项做标记，
// 再次 Snapshot 时成员表 Tags 顺序与内容必须保持原样、未被污染。
func TestSnapshotTagsAliasing(t *testing.T) {
	c := newTestCluster(t)
	mustJoin(t, c, JoinSpec{ID: "n1", Addr: "127.0.0.1:10", Tags: []string{"c", "a", "b"}})
	mustJoin(t, c, JoinSpec{ID: "n2", Addr: "127.0.0.1:11", Tags: []string{"z", "y"}})

	origN1 := mustMember(t, c, "n1")
	origN2 := mustMember(t, c, "n2")

	// 第一次快照：拿到 Tags 后本地排序并改写一项做标记（模拟巡检脚本）。
	snap1 := c.Snapshot()
	for i := range snap1 {
		sort.Strings(snap1[i].Tags) // 本地排序
		snap1[i].Tags[0] = "MARKED:" + snap1[i].ID // 改写一项做标记
	}

	// 再次快照：成员表 Tags 必须保持原始顺序与内容，未被污染。
	snap2 := c.Snapshot()
	gotN1 := findView(snap2, "n1")
	gotN2 := findView(snap2, "n2")
	if !equalTags(gotN1.Tags, origN1.Tags) {
		t.Fatalf("n1 tags 被污染:\nwant %v\n got %v", origN1.Tags, gotN1.Tags)
	}
	if !equalTags(gotN2.Tags, origN2.Tags) {
		t.Fatalf("n2 tags 被污染:\nwant %v\n got %v", origN2.Tags, gotN2.Tags)
	}
}

func TestSnapshotTagsShareNoBackingArray(t *testing.T) {
	c := newTestCluster(t)
	mustJoin(t, c, JoinSpec{ID: "n1", Addr: "127.0.0.1:10", Tags: []string{"a", "b", "c"}})

	s1 := c.Snapshot()
	s2 := c.Snapshot()
	if len(s1) != 1 || len(s2) != 1 {
		t.Fatalf("期望各 1 个成员，got %d / %d", len(s1), len(s2))
	}
	// 写 s1 的返回值，s2 不应受影响（说明不共享底层切片）。
	s1[0].Tags[0] = "MUTATED"
	if s2[0].Tags[0] != "a" {
		t.Fatalf("两次 Snapshot 共享底层切片: s2[0].Tags[0]=%q", s2[0].Tags[0])
	}
}

func mustJoin(t *testing.T, c *Cluster, spec JoinSpec) {
	t.Helper()
	if err := c.Join(spec); err != nil {
		t.Fatalf("Join %s: %v", spec.ID, err)
	}
}

func mustMember(t *testing.T, c *Cluster, id string) MemberView {
	t.Helper()
	v, ok := c.Member(id)
	if !ok {
		t.Fatalf("Member %s 不存在", id)
	}
	return v
}

func findView(list []MemberView, id string) MemberView {
	for _, v := range list {
		if v.ID == id {
			return v
		}
	}
	return MemberView{}
}

func equalTags(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
