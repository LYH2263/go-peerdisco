package peerdisco

import "time"

// State 成员存活状态。
type State int

const (
	StateAlive State = iota
	StateSuspect
	StateDead
	StateLeft
)

func (s State) String() string {
	switch s {
	case StateAlive:
		return "alive"
	case StateSuspect:
		return "suspect"
	case StateDead:
		return "dead"
	case StateLeft:
		return "left"
	default:
		return "unknown"
	}
}

// MemberView 对外成员快照。
type MemberView struct {
	ID       string
	Addr     string
	State    State
	Incarnation uint64
	Meta     map[string]string
	Tags     []string
	JoinedAt time.Time
	Updated  time.Time
}

// JoinSpec Join 入参。
type JoinSpec struct {
	ID   string
	Addr string
	Meta map[string]string
	Tags []string
}

// SuspicionEvent 怀疑队列项。
type SuspicionEvent struct {
	ID        string
	Since     time.Time
	Deadline  time.Time
	From      string
}
