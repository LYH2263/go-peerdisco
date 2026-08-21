package clock

import "time"

type Clock interface {
	Now() time.Time
}

type Real struct{}

func (Real) Now() time.Time { return time.Now() }

// Fake 可拨快进。
type Fake struct {
	t time.Time
}

func NewFake(t time.Time) *Fake { return &Fake{t: t} }

func (f *Fake) Now() time.Time { return f.t }

func (f *Fake) Advance(d time.Duration) { f.t = f.t.Add(d) }

func (f *Fake) Set(t time.Time) { f.t = t }
