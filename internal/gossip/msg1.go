package gossip

// Msg1 辅助消息构造（扩展协议面）。
type Msg1 struct {
	Seq  uint64
	From string
	To   string
	Blob []byte
}

func NewMsg1(seq uint64, from, to string, blob []byte) Msg1 {
	cp := append([]byte(nil), blob...)
	return Msg1{Seq: seq, From: from, To: to, Blob: cp}
}

func (m Msg1) Size() int {
	return 16 + len(m.From) + len(m.To) + len(m.Blob)
}

func (m Msg1) Clone() Msg1 {
	return NewMsg1(m.Seq, m.From, m.To, m.Blob)
}

func ValidateMsg1(m Msg1) bool {
	return m.From != "" && m.To != "" && m.Seq > 0
}
