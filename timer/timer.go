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

func Current() (*Timer, bool) {
	return curr, curr != nil
}

func CanModify(unit wow.UnitID) bool {
	return unitIsRaidAssistant(unit)
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
