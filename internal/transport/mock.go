package transport

import "sync"

// Mock 记录发送。
type Mock struct {
	mu      sync.Mutex
	Sent    []Packet
	Closed  bool
	FailAll bool
}

type Packet struct {
	Addr    string
	Payload []byte
}

func (m *Mock) Send(addr string, payload []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Closed {
		return errClosed
	}
	if m.FailAll {
		return errClosed
	}
	cp := append([]byte(nil), payload...)
	m.Sent = append(m.Sent, Packet{Addr: addr, Payload: cp})
	return nil
}

func (m *Mock) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Closed = true
	return nil
}

func (m *Mock) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Sent)
}
