package transport

import (
	"net"
	"sync"
)

// Transport 发包抽象。
type Transport interface {
	Send(addr string, payload []byte) error
	Close() error
}

// UDPTransport 真实 UDP。
type UDPTransport struct {
	mu   sync.Mutex
	conn *net.UDPConn
}

func ListenUDP(addr string) (*UDPTransport, error) {
	ua, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	c, err := net.ListenUDP("udp", ua)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{conn: c}, nil
}

func (u *UDPTransport) Send(addr string, payload []byte) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	ua, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	_, err = u.conn.WriteToUDP(payload, ua)
	return err
}

func (u *UDPTransport) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.conn == nil {
		return nil
	}
	err := u.conn.Close()
	u.conn = nil
	return err
}

type closedErr string

func (e closedErr) Error() string { return string(e) }

const errClosed = closedErr("transport: closed")
