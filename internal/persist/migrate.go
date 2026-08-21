package persist

import "example.com/peerdisco/internal/member"

// MigrateV1 填充缺省字段。
func MigrateV1(s Snapshot) Snapshot {
	if s.Meta == nil {
		s.Meta = map[string]map[string]string{}
	}
	if s.Members == nil {
		s.Members = []member.Member{}
	}
	for i := range s.Members {
		if s.Members[i].Incarnation == 0 {
			s.Members[i].Incarnation = 1
		}
	}
	return s
}

// FilterAlive 仅保留 Alive。
func FilterAlive(s Snapshot) Snapshot {
	out := s
	kept := make([]member.Member, 0, len(s.Members))
	for _, m := range s.Members {
		if m.State == member.Alive {
			kept = append(kept, m)
		}
	}
	out.Members = kept
	return out
}
