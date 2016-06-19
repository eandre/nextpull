package nextpull

import (
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/boss"
	"github.com/eandre/nextpull/ready"
	"github.com/eandre/nextpull/timer"
	"github.com/eandre/nextpull/ui"
)

func init() {
	b := boss.Find(1788)
	t := &timer.Timer{
		Creator: "foo",
		Started: wow.GetTime(),
		ETA:     wow.GetTime() + 30,
		Boss:    b,
	}
	ui.Show(t, []ready.Item{
		&ready.RunBackItem{Boss: b},
		ready.NewDummyItem("Flask up", "Flask expires in |cffff00005 minutes|r", "Interface\\Icons\\Trade_Alchemy_Dpotion_d11", false),
		&ready.FoodItem{},
		//ready.NewDummyItem("Repair", "Item durability at |cff00ff00100%|r", "Interface\\Icons\\Ability_Repair", true),
	})
}
