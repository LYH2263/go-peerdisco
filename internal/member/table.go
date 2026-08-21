package member

import "fmt"

// Table 成员表。
type Table struct {
	max  int
	byID map[string]Member
	order []string
}

func NewTable(max int) *Table {
	if max <= 0 {
		max = 1024
	}
	return &Table{max: max, byID: make(map[string]Member)}
}

func (t *Table) Add(m Member) error {
	if _, ok := t.byID[m.ID]; ok {
		return fmt.Errorf("member exists")
	}
	if len(t.byID) >= t.max {
		return fmt.Errorf("member table full")
	}
	t.byID[m.ID] = CloneMember(m)
	t.order = append(t.order, m.ID)
	return nil
}

func (t *Table) Get(id string) (Member, bool) {
	m, ok := t.byID[id]
	if !ok {
		return Member{}, false
	}
	return CloneMember(m), true
}

func (t *Table) Update(m Member) error {
	if _, ok := t.byID[m.ID]; !ok {
		return fmt.Errorf("missing")
	}
	t.byID[m.ID] = CloneMember(m)
	return nil
}

func (t *Table) Remove(id string) error {
	if _, ok := t.byID[id]; !ok {
		return fmt.Errorf("missing")
	}
	delete(t.byID, id)
	out := t.order[:0]
	for _, x := range t.order {
		if x != id {
			out = append(out, x)
		}
	}
	t.order = out
	return nil
}

func (t *Table) List() []Member {
	out := make([]Member, 0, len(t.order))
	for _, id := range t.order {
		if m, ok := t.byID[id]; ok {
			out = append(out, CloneMember(m))
		}
	}
	return out
}

func (t *Table) PickHelpers(exclude string, n int) []Member {
	out := make([]Member, 0, n)
	for _, id := range t.order {
		if id == exclude {
			continue
		}
		m := t.byID[id]
		if m.State != Alive {
			continue
		}
		out = append(out, CloneMember(m))
		if len(out) >= n {
			break
		}
	}
	return out
}

func (t *Table) Replace(list []Member) {
	t.byID = make(map[string]Member, len(list))
	t.order = make([]string, 0, len(list))
	for _, m := range list {
		t.byID[m.ID] = CloneMember(m)
		t.order = append(t.order, m.ID)
	}
}

func (t *Table) Clear() {
	t.byID = make(map[string]Member)
	t.order = nil
}

func (t *Table) Len() int { return len(t.byID) }
