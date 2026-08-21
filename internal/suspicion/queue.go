package suspicion

import (
	"time"

	"example.com/peerdisco/internal/clock"
)

type Item struct {
	ID       string
	From     string
	Since    time.Time
	Deadline time.Time
}

type Queue struct {
	timeout time.Duration
	clk     clock.Clock
	items   map[string]Item
}

func NewQueue(timeout time.Duration, clk clock.Clock) *Queue {
	return &Queue{timeout: timeout, clk: clk, items: make(map[string]Item)}
}

func (q *Queue) Add(id, from string, since time.Time) {
	q.items[id] = Item{
		ID:       id,
		From:     from,
		Since:    since,
		Deadline: since.Add(q.timeout),
	}
}

func (q *Queue) Cancel(id string) { delete(q.items, id) }

func (q *Queue) Deadline(id string) (time.Time, bool) {
	it, ok := q.items[id]
	return it.Deadline, ok
}

func (q *Queue) List() []Item {
	out := make([]Item, 0, len(q.items))
	for _, it := range q.items {
		out = append(out, it)
	}
	return out
}

func (q *Queue) Clear() { q.items = make(map[string]Item) }

func (q *Queue) Expired(now time.Time) []string {
	var ids []string
	for id, it := range q.items {
		if !now.Before(it.Deadline) {
			ids = append(ids, id)
		}
	}
	return ids
}
