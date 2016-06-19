package pulldetect

import (
	"github.com/eandre/lunar-wow/pkg/luastrings"
	"github.com/eandre/lunar-wow/pkg/wow"
)

var registry = []func(combat bool){}

func Register(f func(combat bool)) {
	registry = append(registry, f)
}

func trigger(combat bool) {
	for _, f := range registry {
		f(combat)
	}
}

func onEvent(event string, args []interface{}) {
	if event == "ENCOUNTER_START" {
		trigger(true)
	} else if event == "CHAT_MSG_ADDON" {
		prefix := args[0].(string)
		message := args[1].(string)
		if prefix == "BigWigs" && luastrings.HasPrefix(message, "T:BWPull") {
			trigger(false)
		} else if prefix == "D$" && luastrings.HasPrefix(message, "PT\t") {
			trigger(false)
		}
	}
}

func init() {
	wow.RegisterAddonMessagePrefix("BigWigs")
	wow.RegisterAddonMessagePrefix("D4")
	wow.RegisterEvent("CHAT_MSG_ADDON", onEvent)
	wow.RegisterEvent("ENCOUNTER_START", onEvent)
}
