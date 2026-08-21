package gossip

import "hash/fnv"

// Checksum 简易校验。
func Checksum(b []byte) uint32 {
	h := fnv.New32a()
	_, _ = h.Write(b)
	return h.Sum32()
}

// FrameSize 估计帧长。
func FrameSize(from, to string, blob int) int {
	return 1 + 2 + len(from) + 2 + len(to) + blob
}

// SplitAddr host:port 粗分。
func SplitAddr(addr string) (host, port string) {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i], addr[i+1:]
		}
	}
	return addr, ""
}

// MaxPayload 默认上限。
const MaxPayload = 1200

func ClampPayload(b []byte) []byte {
	if len(b) <= MaxPayload {
		return b
	}
	return b[:MaxPayload]
}
