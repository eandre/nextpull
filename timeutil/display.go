package timeutil

import "github.com/eandre/lunar-wow/pkg/luamath"

func Display(dur float32) (mins, secs int, negative bool) {
	neg := false
	if dur < 0 {
		dur = -dur
		dur = float32(luamath.Floor(dur))
		neg = dur > 0
	}

	s := luamath.Ceil(dur)
	m := luamath.Floor(float32(s) / 60)
	s = s % 60
	return m, s, neg
}
