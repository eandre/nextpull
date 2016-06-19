package timer

import (
	"github.com/eandre/lunar-wow/pkg/luastrings"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/boss"
)

const msgPrefix = "NP_TIMER"

func broadcastStartTimer(t *Timer) {
	// No point broadcasting if we cannot start timers
	if !CanModify("player") {
		return
	}

	now := wow.GetTime()
	from := luastrings.ToString(t.Started - now)
	to := luastrings.ToString(t.ETA - now)
	eid := luastrings.ToString(t.Boss.EncounterID)

	msg := "TIMERSTART " + eid + " " + string(t.Creator) + " " + from + " " + to
	wow.SendAddonMessage(msgPrefix, msg, wow.AddonChatTypeRaid, nil)
}

func broadcastStopTimers() {
	if !CanModify("player") {
		return
	}
	wow.SendAddonMessage(msgPrefix, "TIMERSTOP", wow.AddonChatTypeRaid, nil)
}

func init() {
	if !wow.RegisterAddonMessagePrefix(msgPrefix) {
		println("Could not register NextPull timer messages")
	}

	wow.RegisterEvent("CHAT_MSG_ADDON", func(event string, args []interface{}) {
		prefix := args[0].(string)
		message := args[1].(string)
		sender := args[3].(string)
		if prefix != msgPrefix || !nameIsRaidAssistant(sender) {
			return
		}

		if message == "TIMERSTOP" {
			// Stop current timers
			curr = nil
			return
		}

		if !luastrings.HasPrefix(message, "TIMERSTART") {
			// Unrecognized message
			return
		}

		parts := luastrings.Split(" ", message, -1)
		if len(parts) < 5 {
			// Bad format
			return
		}

		// Find the encounter
		eid := boss.EncounterID(luastrings.ToInt(parts[1]))
		b := boss.Find(eid)

		now := wow.GetTime()
		curr = &Timer{
			Creator: wow.GUID(parts[2]),
			Started: now + wow.Time(luastrings.ToFloat(parts[3])),
			ETA:     now + wow.Time(luastrings.ToFloat(parts[4])),
			Boss:    b,
		}
	})
}
