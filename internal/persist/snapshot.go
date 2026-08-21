package persist

import "example.com/peerdisco/internal/member"

// Empty 空快照。
func Empty(selfID, selfAddr string) Snapshot {
	return Snapshot{
		SelfID:   selfID,
		SelfAddr: selfAddr,
		Members:  []member.Member{},
		Meta:     map[string]map[string]string{},
	}
}

// MemberCount 统计。
func (s Snapshot) MemberCount() int { return len(s.Members) }
