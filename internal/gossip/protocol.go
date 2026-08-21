package gossip

import "encoding/binary"

const (
	TypePing byte = 1
	TypeAck  byte = 2
	TypeInd  byte = 3
	TypeMeta byte = 4
)

func EncodePing(from, to string) []byte {
	return encode(TypePing, from, to)
}

func EncodeAck(from, to string) []byte {
	return encode(TypeAck, from, to)
}

func EncodeIndirect(from, to string) []byte {
	return encode(TypeInd, from, to)
}

func encode(t byte, a, b string) []byte {
	ab := []byte(a)
	bb := []byte(b)
	out := make([]byte, 1+2+len(ab)+2+len(bb))
	out[0] = t
	binary.BigEndian.PutUint16(out[1:], uint16(len(ab)))
	copy(out[3:], ab)
	off := 3 + len(ab)
	binary.BigEndian.PutUint16(out[off:], uint16(len(bb)))
	copy(out[off+2:], bb)
	return out
}

func Decode(b []byte) (typ byte, a, c string, err error) {
	if len(b) < 5 {
		return 0, "", "", errShort
	}
	typ = b[0]
	n := int(binary.BigEndian.Uint16(b[1:]))
	if len(b) < 3+n+2 {
		return 0, "", "", errShort
	}
	a = string(b[3 : 3+n])
	off := 3 + n
	m := int(binary.BigEndian.Uint16(b[off:]))
	if len(b) < off+2+m {
		return 0, "", "", errShort
	}
	c = string(b[off+2 : off+2+m])
	return typ, a, c, nil
}

type shortError string

func (e shortError) Error() string { return string(e) }

const errShort = shortError("gossip: short frame")
