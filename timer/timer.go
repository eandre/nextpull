package timer

import (
	"github.com/eandre/lunar-wow/pkg/luastrings"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/boss"
)

type Timer struct {
	Creator wow.GUID
	Started wow.Time
	ETA     wow.Time
	Boss    *boss.Boss
}

var curr *Timer

func newTimer(t *Timer) {
	curr = t
	for _, f := range registry {
		f(t, false)
	}
}

func stopTimer() {
	if curr == nil {
		return
	}

	t := curr
	curr = nil
	for _, f := range registry {
		f(t, true)
	}
}

func CanModify(unit wow.UnitID) bool {
	return unitIsRaidAssistant(unit)
}

func StopAll() {
	broadcastStopTimers()
}

func Start(t *Timer) {
	broadcastStartTimer(t)
}

var registry []func(*Timer, bool)

func RegisterCallback(f func(t *Timer, stopped bool)) {
	registry = append(registry, f)
}

func unitIsRaidAssistant(unit wow.UnitID) bool {
	for i := wow.GroupIndex(1); i <= wow.GetNumGroupMembers(); i++ {
		uid := wow.UnitID("raid" + luastrings.ToString(i))
		if wow.UnitIsUnit(unit, uid) {
			_, rank, _, _, _, _, _, _, _, _, _ := wow.GetRaidRosterInfo(i)
			return rank >= wow.RaidRankAssistant
		}
	}
	return false
}

func nameIsRaidAssistant(name string) bool {
	for i := wow.GroupIndex(1); i <= wow.GetNumGroupMembers(); i++ {
		uid := wow.UnitID("raid" + luastrings.ToString(i))
		raidName, realm := wow.UnitName(uid)
		if name == raidName || name == (raidName+"-"+realm) {
			_, rank, _, _, _, _, _, _, _, _, _ := wow.GetRaidRosterInfo(i)
			return rank >= wow.RaidRankAssistant
		}
	}
	return false
}
