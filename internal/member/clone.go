package member

// CloneTags 深拷贝标签切片，避免与内部存储共享底层数组。
// 成员快照必须返回独立副本，否则调用方就地改 Tags[0] 会写穿到
// 集群成员表，进而被 gossip/persist 带出脏标签。
func CloneTags(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// CloneMember 深拷贝成员。
func CloneMember(m Member) Member {
	out := m
	out.Tags = CloneTags(m.Tags)
	return out
}

// CloneList 深拷贝成员列表。
func CloneList(src []Member) []Member {
	if src == nil {
		return nil
	}
	dst := make([]Member, len(src))
	for i := range src {
		dst[i] = CloneMember(src[i])
	}
	return dst
}
