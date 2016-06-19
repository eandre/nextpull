package nextpull

import (
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/ready"
	"github.com/eandre/nextpull/timer"
	"github.com/eandre/nextpull/ui"
)

type DB struct {
	Items []ready.Item
}

var DefaultItems = &DB{
	Items: []ready.Item{
		&ready.FoodItem{},
	},
}

func init() {
	t := &timer.Timer{
		Creator: "foo",
		Started: wow.GetTime(),
		ETA:     wow.GetTime() + 30,
	}
	ui.Show(t, DefaultItems.Items)
}
