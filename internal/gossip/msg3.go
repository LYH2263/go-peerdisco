package gossip

// Msg3 辅助消息构造（扩展协议面）。
type Msg3 struct {
	Seq  uint64
	From string
	To   string
	Blob []byte
}

func NewMsg3(seq uint64, from, to string, blob []byte) Msg3 {
	cp := append([]byte(nil), blob...)
	return Msg3{Seq: seq, From: from, To: to, Blob: cp}
}

func (m Msg3) Size() int {
	return 16 + len(m.From) + len(m.To) + len(m.Blob)
}

func (m Msg3) Clone() Msg3 {
	return NewMsg3(m.Seq, m.From, m.To, m.Blob)
}

func ValidateMsg3(m Msg3) bool {
	return m.From != "" && m.To != "" && m.Seq > 0
}
