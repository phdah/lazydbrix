package tui

import (
	"github.com/rivo/tview"
)

type Text struct {
	Text *tview.TextView
}

func NewText() *Text {
	detailText := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true)

	return &Text{detailText}
}

func (dt *Text) Setup(name string) {
	dt.Text.SetBorder(true).SetTitle(name)
}
