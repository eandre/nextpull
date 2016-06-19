package ui

import (
	"github.com/eandre/lunar-wow/pkg/luamath"
	"github.com/eandre/lunar-wow/pkg/widget"
	"github.com/eandre/lunar-wow/pkg/wow"
	"github.com/eandre/nextpull/ready"
	"github.com/eandre/nextpull/timer"
)

type readyFrame struct {
	root    widget.Frame
	timer   widget.FontString
	items   []*itemFrame
	botLine widget.Texture
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

func updateItemFrame(f *itemFrame, item ready.Item) {
	f.ri = item
	f.name.SetText(item.Name())
	f.subText.SetText(item.Description())
	f.icon.SetTexture(item.Icon())

	if item.Ready("player") {
		f.icon.SetAlpha(0.5)
		f.name.SetAlpha(0.5)
		f.subText.SetAlpha(0.5)
		f.check.Show()
	} else {
		f.icon.SetAlpha(1)
		f.name.SetAlpha(1)
		f.subText.SetAlpha(1)
		f.check.Hide()
	}
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

	for _, f := range frame.items {
		updateItemFrame(f, f.ri)
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
	sign := ""
	if dur < 0 {
		dur = -dur
		sign = "+"
	}

	secs := luamath.Ceil(float32(dur))
	mins := luamath.Floor(float32(secs) / 60)
	secs = secs % 60
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

	topGlow := root.CreateTexture()
	topGlow.SetPoint("BOTTOM", root, "TOP", 0, 7)
	topGlow.SetTexture("Interface\\LevelUp\\LevelUpTex")
	topGlow.SetSize(226, 117)
	topGlow.SetTexCoord(0.55859375, 1, 0.240234375, 0.466796875)

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
		botLine: botLine,
		timer:   timerText,
	}
}
