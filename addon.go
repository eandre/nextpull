package nextpull

import (
	"github.com/eandre/lunar-shim/ace/aceconsole"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/boss"
	"github.com/eandre/nextpull/pulldetect"
	"github.com/eandre/nextpull/ready"
	"github.com/eandre/nextpull/timer"
	"github.com/eandre/nextpull/ui"
)

var currItems []ready.Item
var currTimer *timer.Timer

func stopTimer() {
	currItems = nil
	currTimer = nil
	ui.Hide()
}

func startTimer(t *timer.Timer, items []ready.Item) {
	currTimer = t
	currItems = items
	ui.Show(t, items)
}

func init() {
	timer.RegisterCallback(func(t *timer.Timer, stopped bool) {
		if stopped {
			stopTimer()
			return
		}

		startTimer(t, []ready.Item{
			&ready.RunBackItem{Boss: t.Boss},
			ready.NewFlaskItem(t.Boss),
			ready.NewWellFedItem(t.Boss),
		})
	})

	pulldetect.Register(func(combat bool) {
		ui.Hide()
	})

	b := boss.Find(1788)
	t := &timer.Timer{
		Creator: "foo",
		Started: wow.GetTime(),
		ETA:     wow.GetTime() + 5,
		Boss:    b,
	}
	startTimer(t, []ready.Item{
		&ready.RunBackItem{Boss: b},
		ready.NewFlaskItem(b),
		ready.NewWellFedItem(b),
	})

	aceconsole.RegisterChatCommand("nextpull", slashCmd)
	aceconsole.RegisterChatCommand("np", slashCmd)
}

var currBoss *boss.Boss

func slashCmd(msg string) {
	if msg == "boss" {
		b := boss.FindClosest()
		if b == nil {
			println("Could not find a boss in this instance.")
			return
		}
		println("Setting boss to: " + b.Name + ". Use /np to start a pull timer.")
		currBoss = b
	} else if msg == "stop" {
		timer.StopAll()
		stopTimer()
	} else if msg == "" || msg == "start" {
		if !timer.CanModify("player") {
			println("Need to be raid leader or assistant to be able to start pull timer.")
			return
		}
		if currTimer != nil {
			println("There is already a timer going. To create a new one, first stop it with /np stop")
			return
		}
		if currBoss == nil {
			currBoss = boss.FindClosest()
			if currBoss == nil {
				println("No boss found. Enter a supported raid instance?")
				return
			}
		}

		t := &timer.Timer{
			Creator: wow.UnitGUID("player"),
			Started: wow.GetTime(),
			ETA:     wow.GetTime() + (60 * 5),
			Boss:    currBoss,
		}
		timer.Start(t)
	}
}
