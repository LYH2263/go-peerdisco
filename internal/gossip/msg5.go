package gossip

// Msg5 辅助消息构造（扩展协议面）。
type Msg5 struct {
	Seq  uint64
	From string
	To   string
	Blob []byte
}

func NewMsg5(seq uint64, from, to string, blob []byte) Msg5 {
	cp := append([]byte(nil), blob...)
	return Msg5{Seq: seq, From: from, To: to, Blob: cp}
}

func (m Msg5) Size() int {
	return 16 + len(m.From) + len(m.To) + len(m.Blob)
}

func (m Msg5) Clone() Msg5 {
	return NewMsg5(m.Seq, m.From, m.To, m.Blob)
}

func ValidateMsg5(m Msg5) bool {
	return m.From != "" && m.To != "" && m.Seq > 0
}
