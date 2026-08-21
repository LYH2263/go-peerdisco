package gossip

// Msg4 辅助消息构造（扩展协议面）。
type Msg4 struct {
	Seq  uint64
	From string
	To   string
	Blob []byte
}

func NewMsg4(seq uint64, from, to string, blob []byte) Msg4 {
	cp := append([]byte(nil), blob...)
	return Msg4{Seq: seq, From: from, To: to, Blob: cp}
}

func (m Msg4) Size() int {
	return 16 + len(m.From) + len(m.To) + len(m.Blob)
}

func (m Msg4) Clone() Msg4 {
	return NewMsg4(m.Seq, m.From, m.To, m.Blob)
}

func ValidateMsg4(m Msg4) bool {
	return m.From != "" && m.To != "" && m.Seq > 0
}
