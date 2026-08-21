package gossip

import "bytes"

// EncodeMeta 元数据 gossip 载荷。
func EncodeMeta(id string, blob []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(TypeMeta)
	buf.WriteByte(byte(len(id)))
	buf.WriteString(id)
	buf.Write(blob)
	return buf.Bytes()
}

// Pad 填充到最小帧（测试辅助）。
func Pad(b []byte, n int) []byte {
	if len(b) >= n {
		return b
	}
	out := make([]byte, n)
	copy(out, b)
	return out
}
