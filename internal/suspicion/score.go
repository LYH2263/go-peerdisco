package suspicion

import "time"

// Score 简单怀疑分：超时越近分越高。
func Score(it Item, now time.Time) float64 {
	total := it.Deadline.Sub(it.Since).Seconds()
	if total <= 0 {
		return 1
	}
	elapsed := now.Sub(it.Since).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}
	s := elapsed / total
	if s > 1 {
		return 1
	}
	return s
}

// SortByDeadline 按截止排序（插入排序，成员少）。
func SortByDeadline(items []Item) []Item {
	out := append([]Item(nil), items...)
	for i := 1; i < len(out); i++ {
		j := i
		for j > 0 && out[j].Deadline.Before(out[j-1].Deadline) {
			out[j], out[j-1] = out[j-1], out[j]
			j--
		}
	}
	return out
}
