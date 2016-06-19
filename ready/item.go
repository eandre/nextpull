package ready

import (
	"github.com/eandre/lunar-shim/hbd"
	"github.com/eandre/lunar-wow/pkg/luastrings"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/boss"
	"github.com/eandre/nextpull/timeutil"
)

type Item interface {
	Name() string
	Description() string
	Icon() string
	Ready(unit wow.UnitID) bool
}

type RunBackItem struct {
	Boss *boss.Boss
}

func (i *RunBackItem) Name() string {
	return "Run back"
}

func (i *RunBackItem) Description() string {
	return "Get to |cffffffff" + i.Boss.Name + "|r"
}

func (i *RunBackItem) Icon() string {
	return "Interface\\Icons\\Ability_Hunter_MarkedForDeath"
}

func (i *RunBackItem) Ready(unit wow.UnitID) bool {
	x, y, inst := hbd.UnitWorldPosition(unit)
	if inst != i.Boss.Raid.InstanceID {
		return false
	}
	dist := hbd.WorldDistance(inst, x, y, i.Boss.X, i.Boss.Y)
	return dist <= 30
}

type AuraItem struct {
	boss        *boss.Boss
	name        string
	icon        string
	desc        string
	missingDesc string
	ids         map[int64]bool
}

func NewWellFedItem(b *boss.Boss) *AuraItem {
	return &AuraItem{
		boss:        b,
		name:        "Eat Food",
		missingDesc: "No food buff",
		icon:        "Interface\\Icons\\Spell_Misc_Food",
		ids:         wellFedMap,
	}
}

func NewFlaskItem(b *boss.Boss) *AuraItem {
	return &AuraItem{
		boss:        b,
		name:        "Flask up",
		missingDesc: "No flask buff",
		icon:        "Interface\\Icons\\Trade_Alchemy_Dpotion_d11",
		ids:         flaskMap,
	}
}

func (i *AuraItem) Name() string {
	return i.name
}

func (i *AuraItem) Icon() string {
	return i.icon
}

func (i *AuraItem) Description() string {
	return i.desc
}

func (i *AuraItem) Ready(unit wow.UnitID) bool {
	dur, ok := i.getDur(unit)
	if wow.UnitIsUnit(unit, "player") {
		i.updateDesc(dur)
	}
	return ok && dur > i.boss.FightDuration()
}

func (i *AuraItem) getDur(unit wow.UnitID) (float32, bool) {
	for idx := 1; ; idx++ {
		_, _, _, _, _, _, expires, _, _, _, spellID, _, _, _, _, _ := wow.UnitAura(unit, idx, "HELPFUL")
		if spellID == 0 {
			return -1, false
		} else if i.ids[spellID] {
			dt := float32(expires - wow.GetTime())
			return dt, true
		}
	}
}

func (i *AuraItem) updateDesc(dur float32) {
	mins, secs, _ := timeutil.Display(dur)
	if dur <= 0 {
		i.desc = "|cffff0000" + i.missingDesc + "|r"
	} else if dur > i.boss.FightDuration() {
		i.desc = "Expires in |cffffffff" + luastrings.ToString(mins+1) + " minutes|r"
	} else if mins > 0 {
		i.desc = "Expires in |cffff0000" + luastrings.ToString(mins+1) + " minutes|r"
	} else {
		i.desc = "Expires in |cffff0000" + luastrings.ToString(secs) + " seconds|r"
	}
}

var wellFedMap = map[int64]bool{
	180749: true,
	180745: true,
	180750: true,
	180748: true,
	180746: true,
	180747: true,
	188534: true,
}

var flaskMap = map[int64]bool{
	156079: true,
}

type DummyItem struct {
	name        string
	description string
	icon        string
	ready       bool
}

func (i *DummyItem) Name() string               { return i.name }
func (i *DummyItem) Description() string        { return i.description }
func (i *DummyItem) Icon() string               { return i.icon }
func (i *DummyItem) Ready(unit wow.UnitID) bool { return i.ready }

func NewDummyItem(name, description, icon string, ready bool) *DummyItem {
	return &DummyItem{
		name:        name,
		description: description,
		icon:        icon,
		ready:       ready,
	}
}
