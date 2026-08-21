package peerdisco

import "errors"

var (
	ErrClosed      = errors.New("peerdisco: closed")
	ErrExists      = errors.New("peerdisco: already joined")
	ErrNotMember   = errors.New("peerdisco: not a member")
	ErrInvalid     = errors.New("peerdisco: invalid argument")
	ErrNoTransport = errors.New("peerdisco: no transport")
	ErrPersist     = errors.New("peerdisco: persist failed")
	ErrCanceled    = errors.New("peerdisco: canceled")
	ErrDead        = errors.New("peerdisco: member dead")
	ErrSuspect     = errors.New("peerdisco: member suspected")
	ErrTimeout     = errors.New("peerdisco: timeout")
)
