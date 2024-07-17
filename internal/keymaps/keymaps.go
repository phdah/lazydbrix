package keymaps

import (
	"github.com/gdamore/tcell/v2"
	"github.com/phdah/lazydbrix/internal/utils"
	"github.com/rivo/tview"
)

// Set keymaps for a tview.List
func SetListKeymaps(list *tview.List) {
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				utils.MoveListDown(list)
				return nil
			case 'k':
				utils.MoveListUp(list)
				return nil
			}
        case tcell.KeyEnter:
            utils.MakeListSelection(list)
            return nil
		}
		return event
	})
}

// Set keymaps for a tview.Flex
func SetFlexKeymaps(app *tview.Application, flex *tview.Flex) {
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h':
				utils.MoveFlexItemUp(app, flex)
				return nil
			case 'l':
				utils.MoveFlexItemDown(app, flex)
				return nil
			}
		}
		return event
	})
}

// Set keymaps for the main tview.Flex
func SetMainFlexKeymaps(app *tview.Application, flex *tview.Flex) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			utils.MoveFlexRight(app, flex)
			return nil
		case tcell.KeyBacktab:
			utils.MoveFlexLeft(app, flex)
			return nil
		}
		return event
	})
}
