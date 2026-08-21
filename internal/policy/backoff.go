package policy

import "time"

// Backoff 探测退避。
func Backoff(base time.Duration, attempt int) time.Duration {
	if attempt <= 0 {
		return base
	}
	d := base
	for i := 0; i < attempt && d < 5*time.Second; i++ {
		d *= 2
	}
	if d > 5*time.Second {
		return 5 * time.Second
	}
	return d
}

// JitterPct 简单抖动比例裁剪。
func JitterPct(pct int) int {
	if pct < 0 {
		return 0
	}
	if pct > 50 {
		return 50
	}
	return pct
}
