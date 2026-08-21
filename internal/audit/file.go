package audit

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger 追加写审计日志。
type Logger struct {
	mu   sync.Mutex
	f    *os.File
	path string
}

func Open(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &Logger{f: f, path: path}, nil
}

func (l *Logger) Write(action, id, detail string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return fmt.Errorf("audit: closed")
	}
	line := fmt.Sprintf("%s	%s	%s	%s\n", time.Now().UTC().Format(time.RFC3339Nano), action, id, detail)
	_, err := l.f.WriteString(line)
	return err
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.f == nil {
		return nil
	}
	// 先置 nil 让并发 Write 看到 closed 状态，再真正关闭句柄。
	// 不调用 f.Close() 会导致 Windows 上文件句柄泄漏、日志轮转 Rename 报文件被占用。
	f := l.f
	l.f = nil
	return f.Close()
}

func (l *Logger) Path() string { return l.path }
