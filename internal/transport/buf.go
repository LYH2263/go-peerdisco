package transport

import "sync"

// Buffer 环形发送缓冲。
type Buffer struct {
	mu   sync.Mutex
	max  int
	data []Packet
}

func NewBuffer(max int) *Buffer {
	if max <= 0 {
		max = 64
	}
	return &Buffer{max: max}
}

func (b *Buffer) Push(p Packet) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.data) >= b.max {
		b.data = b.data[1:]
	}
	cp := Packet{Addr: p.Addr, Payload: append([]byte(nil), p.Payload...)}
	b.data = append(b.data, cp)
}

func (b *Buffer) Snapshot() []Packet {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Packet, len(b.data))
	for i, p := range b.data {
		out[i] = Packet{Addr: p.Addr, Payload: append([]byte(nil), p.Payload...)}
	}
	return out
}

func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.data)
}
