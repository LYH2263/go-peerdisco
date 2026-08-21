package gossip

// Msg2 辅助消息构造（扩展协议面）。
type Msg2 struct {
	Seq  uint64
	From string
	To   string
	Blob []byte
}

func NewMsg2(seq uint64, from, to string, blob []byte) Msg2 {
	cp := append([]byte(nil), blob...)
	return Msg2{Seq: seq, From: from, To: to, Blob: cp}
}

func (m Msg2) Size() int {
	return 16 + len(m.From) + len(m.To) + len(m.Blob)
}

func (m Msg2) Clone() Msg2 {
	return NewMsg2(m.Seq, m.From, m.To, m.Blob)
}

func ValidateMsg2(m Msg2) bool {
	return m.From != "" && m.To != "" && m.Seq > 0
}
