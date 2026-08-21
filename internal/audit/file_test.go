package audit

import (
	"path/filepath"
	"testing"
)

// TestCloseReleasesHandle 防回归：Close 必须真正关闭文件句柄，
// 否则 Windows 上运维轮转 Rename(audit.log, audit.log.1) 会报文件被占用。
func TestCloseReleasesHandle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	lg, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := lg.Write("join", "id1", "addr1"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := lg.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// 句柄释放后，Rename 应能成功——这是轮转的前提。
	rotated := filepath.Join(dir, "audit.log.1")
	if err := Rotate(path, rotated); err != nil {
		t.Fatalf("Rotate after Close: %v", err)
	}

	// 关闭后再次 Close 不应 panic、不应报错。
	if err := lg.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}

	// 关闭后写应返回 closed 错误，而不是写入已关闭句柄。
	if err := lg.Write("join", "id2", "addr2"); err == nil {
		t.Fatal("Write after Close: expected error, got nil")
	}
}
