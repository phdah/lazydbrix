package colors

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetCustomListColor(list *tview.List) {
	list.SetSelectedBackgroundColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorBlack)
}
