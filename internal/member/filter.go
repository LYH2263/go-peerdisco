package member

// FilterByState 过滤状态。
func FilterByState(list []Member, st State) []Member {
	out := make([]Member, 0, len(list))
	for _, m := range list {
		if m.State == st {
			out = append(out, CloneMember(m))
		}
	}
	return out
}

// IDs 提取 ID。
func IDs(list []Member) []string {
	out := make([]string, len(list))
	for i, m := range list {
		out[i] = m.ID
	}
	return out
}

// FindAddr 查地址。
func FindAddr(list []Member, id string) (string, bool) {
	for _, m := range list {
		if m.ID == id {
			return m.Addr, true
		}
	}
	return "", false
}
