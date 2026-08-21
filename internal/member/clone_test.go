package member

import "testing"

// TestCloneTagsIndependent 断言 CloneTags 返回独立副本：修改返回切片
// 不得写穿源切片。这是管理页把 Tags[0] 改成 mask 却污染到集群成员表、
// 进而通过 gossip 带出脏标签的回归用例。
func TestCloneTagsIndependent(t *testing.T) {
	src := []string{"prod", "cache"}
	clone := CloneTags(src)

	if len(src) > 0 && &src[0] == &clone[0] {
		t.Fatalf("CloneTags 返回的切片与源共享底层数组")
	}

	// 模拟前端脱敏就地改 Tags[0]。
	clone[0] = "mask"

	if src[0] != "prod" {
		t.Fatalf("修改副本写穿了源切片：got %q, want %q", src[0], "prod")
	}
}

func TestCloneTagsNil(t *testing.T) {
	if got := CloneTags(nil); got != nil {
		t.Fatalf("CloneTags(nil) = %v, want nil", got)
	}
}

func TestCloneTagsEmpty(t *testing.T) {
	got := CloneTags([]string{})
	if len(got) != 0 {
		t.Fatalf("CloneTags(empty) len = %d, want 0", len(got))
	}
}

// TestCloneMemberTagsIndependent 断言整条克隆链对 Tags 独立：
// List() 返回的成员快照就地改 Tags 不得影响 Table 内部存储。
func TestCloneMemberTagsIndependent(t *testing.T) {
	tbl := NewTable(8)
	_ = tbl.Add(Member{
		ID:    "n1",
		Addr:  "127.0.0.1:1",
		State: Alive,
		Tags:  []string{"prod", "cache"},
	})

	listed := tbl.List()
	if len(listed) != 1 {
		t.Fatalf("List len = %d, want 1", len(listed))
	}

	// 前端脱敏：把快照里的首标签改成 mask。
	listed[0].Tags[0] = "mask"

	// 重新读取内部存储，真实标签不应被污染。
	got, ok := tbl.Get("n1")
	if !ok {
		t.Fatal("Get(n1) 未找到")
	}
	if got.Tags[0] != "prod" {
		t.Fatalf("快照写穿了内部表：Tags[0] = %q, want %q", got.Tags[0], "prod")
	}
}
