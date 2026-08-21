package member

import "time"

type State int

const (
	Alive State = iota
	Suspect
	Dead
	Left
)

type Member struct {
	ID          string
	Addr        string
	State       State
	Incarnation uint64
	Tags        []string
	JoinedAt    time.Time
	Updated     time.Time
}
