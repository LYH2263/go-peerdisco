package audit

import (
	"os"
	"path/filepath"
)

// Rotate 关闭当前文件并换新路径（调用方须已 Close）。
func Rotate(oldPath, newPath string) error {
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(oldPath); err == nil {
		return os.Rename(oldPath, newPath)
	}
	return nil
}
