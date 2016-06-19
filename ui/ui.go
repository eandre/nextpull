package ui

import (
	"github.com/eandre/lunar-wow/pkg/widget"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/ready"
	"github.com/eandre/nextpull/timer"
	"github.com/eandre/nextpull/timeutil"
)

type readyFrame struct {
	root    widget.Frame
	timer   widget.FontString
	items   []*itemFrame
	botLine widget.Texture
	glow    widget.Texture
	gonged  bool
}

type itemFrame struct {
	root    widget.Frame
	icon    widget.Texture
	name    widget.FontString
	subText widget.FontString
	check   widget.Texture
	ri      ready.Item
}

var frame *readyFrame
var itemFramePool []*itemFrame

var pullTimer *timer.Timer
var readyItems []ready.Item

func Show(t *timer.Timer, items []ready.Item) {
	frame.root.Show()
	frame.gonged = false
	Update(t, items)
}

func Hide() {
	pullTimer = nil
	readyItems = []ready.Item{}
	frame.root.Hide()
}

func Update(t *timer.Timer, items []ready.Item) {
	pullTimer = t
	readyItems = items

	for _, f := range itemFramePool {
		f.root.Hide()
	}
	frame.items = []*itemFrame{}

	var prev *itemFrame
	for i, item := range items {
		f := getItemFrame(i)
		updateItemFrame(f, item)
		frame.items = append(frame.items, f)

		f.root.Show()
		if prev == nil {
			f.root.SetPoint("TOP", frame.root, "TOP", 0, 0)
		} else {
			f.root.SetPoint("TOP", prev.root, "BOTTOM", 0, 0)
		}
		prev = f
	}

	if prev == nil {
		// Nothing to display
		frame.root.SetHeight(1)
		frame.botLine.Hide()
	} else {
		// Encompass all the items
		frame.root.SetHeight(float32(len(items)) * prev.root.GetHeight())
		frame.botLine.Show()
	}
}

func updateItemFrame(f *itemFrame, item ready.Item) bool {
	f.ri = item
	f.name.SetText(item.Name())
	f.subText.SetText(item.Description())
	f.icon.SetTexture(item.Icon())

	if item.Ready("player") {
		f.icon.SetAlpha(0.5)
		f.name.SetAlpha(0.5)
		f.subText.SetAlpha(0.5)
		f.check.Show()
		return true
	}

	f.icon.SetAlpha(1)
	f.name.SetAlpha(1)
	f.subText.SetAlpha(1)
	f.check.Hide()
	return false
}

var itemAcc float32

func itemUpdate(dt float32) {
	if pullTimer == nil {
		return
	}

	// Only update once every 0.5s
	itemAcc += dt
	if itemAcc < 0.5 {
		return
	}
	itemAcc -= 0.5

	allReady := true
	for _, f := range frame.items {
		rdy := updateItemFrame(f, f.ri)
		allReady = allReady && rdy
	}

	shouldGong := false
	if allReady {
		frame.glow.SetVertexColor(0, 1, 0, 1)
	} else if pullTimer.ETA < wow.GetTime() {
		// time exceeded
		frame.glow.SetVertexColor(1, 0, 0, 1)
		shouldGong = true
	} else {
		frame.glow.SetVertexColor(1, 1, 1, 1)
	}

	if shouldGong && !frame.gonged {
		// Gong!
		wow.PlaySound("FX_Scene_AlittlePatience_SiteAttack_Alert", "Master")
		frame.gonged = true
	}
}

var timerAcc float32

func timerUpdate(dt float32) {
	if pullTimer == nil {
		return
	}

	// Only update once every 0.1s
	timerAcc += dt
	if timerAcc < 0.1 {
		return
	}
	timerAcc -= 0.1

	dur := pullTimer.ETA - wow.GetTime()
	mins, secs, neg := timeutil.Display(float32(dur))
	sign := ""
	if neg {
		sign = "+"
	}
	frame.timer.SetFormattedText("%s%02d:%02d", sign, mins, secs)
}

func getItemFrame(idx int) *itemFrame {
	if (idx + 1) <= len(itemFramePool) {
		return itemFramePool[idx]
	}

	root := widget.CreateFrame(frame.root)
	root.SetSize(300, 52)

	icon := root.CreateTexture()
	icon.SetSize(36, 36)
	icon.SetPoint("TOPLEFT", root, "TOPLEFT", 8, -8)

	name := root.CreateFontString()
	name.SetFontObject("GameFontNormalLarge")
	name.SetJustifyH(widget.JustifyLeft)
	name.SetPoint("TOPLEFT", icon, "TOPRIGHT", 10, -6)
	name.SetTextColor(1, 0.82, 0, 1)

	subText := root.CreateFontString()
	subText.SetFontObject("GameFontNormal")
	subText.SetJustifyH(widget.JustifyLeft)
	subText.SetPoint("TOPLEFT", name, "BOTTOMLEFT", 0, 1)
	subText.SetTextColor(0.7, 0.7, 0.7, 1)

	check := root.CreateTexture()
	check.SetDrawLayer(widget.LayerOverlay, 0)
	check.SetSize(24, 24)
	check.SetPoint("RIGHT", root, "right", 0, 0)
	check.SetTexture("Interface\\RaidFrame\\ReadyCheck-Ready")
	check.Hide()

	f := &itemFrame{
		root:    root,
		icon:    icon,
		name:    name,
		subText: subText,
		check:   check,
	}
	itemFramePool = append(itemFramePool, f)
	return f
}

func init() {
	wow.RegisterUpdate(itemUpdate)
	wow.RegisterUpdate(timerUpdate)

	root := widget.CreateFrame(widget.UIParent())
	root.SetSize(418, 72)
	root.SetPoint("TOP", widget.UIParent(), "TOP", 0, -190)
	root.Hide()

	topLine := root.CreateTexture()
	topLine.SetPoint("BOTTOM", root, "TOP", 0, 7)
	topLine.SetTexture("Interface\\LevelUp\\LevelUpTex")
	topLine.SetSize(418, 7)
	topLine.SetTexCoord(0.00195313, 0.81835938, 0.01953125, 0.03320313)

	glow := root.CreateTexture()
	glow.SetPoint("BOTTOM", root, "TOP", 0, 7)
	glow.SetTexture("Interface\\LevelUp\\LevelUpTex")
	glow.SetSize(226, 117)
	glow.SetTexCoord(0.55859375, 1, 0.240234375, 0.466796875)

	botLine := root.CreateTexture()
	botLine.SetPoint("TOP", root, "BOTTOM", 0, 0)
	botLine.SetTexture("Interface\\LevelUp\\LevelUpTex")
	botLine.SetSize(418, 7)
	botLine.SetTexCoord(0.00195313, 0.81835938, 0.01953125, 0.03320313)

	timerText := root.CreateFontString()
	timerText.SetFontObject("GameFont_Gigantic")
	timerText.SetPoint("BOTTOM", root, "TOP", 0, 7)
	timerText.SetDrawLayer(widget.LayerArtwork, 0)
	timerText.SetTextColor(1, 1, 1, 1)

	timerDesc := root.CreateFontString()
	timerDesc.SetFontObject("GameFontNormal")
	timerDesc.SetPoint("BOTTOM", timerText, "TOP", 0, 5)
	timerDesc.SetText("Next Pull")
	timerDesc.SetDrawLayer(widget.LayerArtwork, 0)

	frame = &readyFrame{
		root:    root,
		glow:    glow,
		botLine: botLine,
		timer:   timerText,
	}
}
