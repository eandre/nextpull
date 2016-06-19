package ready

import "github.com/eandre/lunar-wow/pkg/wow"

type Item interface {
	Name() string
	Description() string
	Icon() string
	Ready(unit wow.UnitID) bool
}

type FoodItem struct{}

func (i *FoodItem) Name() string {
	return "Eat Food"
}

func (i *FoodItem) Icon() string {
	return "Interface\\Icons\\Spell_Misc_Food"
}

func (i *FoodItem) Description() string {
	return "Become Well Fed"
}

func (i *FoodItem) Ready(unit wow.UnitID) bool {
	return false
}
