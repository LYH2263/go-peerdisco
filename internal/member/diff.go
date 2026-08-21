package member

// DiffIDs 返回 only-in-a / only-in-b。
func DiffIDs(a, b []Member) (onlyA, onlyB []string) {
	ma := map[string]struct{}{}
	mb := map[string]struct{}{}
	for _, m := range a {
		ma[m.ID] = struct{}{}
	}
	for _, m := range b {
		mb[m.ID] = struct{}{}
	}
	for id := range ma {
		if _, ok := mb[id]; !ok {
			onlyA = append(onlyA, id)
		}
	}
	for id := range mb {
		if _, ok := ma[id]; !ok {
			onlyB = append(onlyB, id)
		}
	}
	return onlyA, onlyB
}

// MergeTags 合并去重。
func MergeTags(a, b []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(a)+len(b))
	for _, x := range append(append([]string{}, a...), b...) {
		if _, ok := seen[x]; ok {
			continue
		}
		seen[x] = struct{}{}
		out = append(out, x)
	}
	return out
}

