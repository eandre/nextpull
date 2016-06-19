package ready

import (
	"github.com/eandre/lunar-shim/hbd"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/boss"
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

type FoodItem struct{}

func (i *FoodItem) Name() string {
	return "Eat Food"
}

func (i *FoodItem) Icon() string {
	return "Interface\\Icons\\Spell_Misc_Food"
}

func (i *FoodItem) Description() string {
	return "Become |cffffffffWell Fed|r"
}

func (i *FoodItem) Ready(unit wow.UnitID) bool {
	for idx := 1; ; idx++ {
		_, _, _, _, _, _, _, _, _, _, spellID, _, _, _, _, _ := wow.UnitAura(unit, idx, "HELPFUL")
		if spellID == 0 {
			return false
		} else if wellFedMap[spellID] {
			return true
		}
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
