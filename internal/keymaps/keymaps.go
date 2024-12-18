package keymaps

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/phdah/lazydbrix/internal/tui"
	"github.com/phdah/lazydbrix/internal/utils"
	"github.com/rivo/tview"
)

// Set keymaps for a tview.List
func SetEnvKeymaps(app *tview.Application, envList *tview.List) {
	originalCapture := envList.GetInputCapture()
	envList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				utils.MoveListDown(envList)
				return nil
			case 'k':
				utils.MoveListUp(envList)
				return nil
			case 'q':
				log.Printf("Trying to quite")
				app.Stop()
				return nil
			}
		}
		// Pass the event to the original handler if not handled here
		if originalCapture != nil {
			return originalCapture(event)
		}
		return event
	})
}

// Helper function to make selections in a list
func MakeListSelection(envList *tview.List, clusterList *tview.List, clusterSelection *tui.ClusterSelection) {
	index := clusterList.GetCurrentItem()
	itemFirstText, itemSecondText := clusterList.GetItemText(index)

	envMainText, _ := envList.GetItemText(envList.GetCurrentItem())
	clusterMainText, clusterSecondaryText := clusterList.GetItemText(clusterList.GetCurrentItem())
	clusterSelection.Profile = envMainText
	clusterSelection.ClusterName = clusterMainText
	clusterSelection.ClusterID = clusterSecondaryText

	coloredItemFirstText := fmt.Sprintf("[green]%s", itemFirstText)
	clusterList.SetItemText(index, coloredItemFirstText, itemSecondText)
}

// Set keymaps for a tview.List
func SetClusterKeymaps(app *tview.Application, envList *tview.List, clusterList *tview.List, clusterSelection *tui.ClusterSelection) {
	clusterList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				utils.MoveListDown(clusterList)
				return nil
			case 'k':
				utils.MoveListUp(clusterList)
				return nil
			case 'q':
				log.Printf("Trying to quite")
				app.Stop()
				return nil
			}
		case tcell.KeyEnter:
			MakeListSelection(envList, clusterList, clusterSelection)
			log.Printf("ClusterSelection has been updated to: %s", clusterSelection.ClusterName)
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
			case 'q':
				log.Printf("Trying to quite")
				app.Stop()
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
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				log.Printf("Trying to quite")
				app.Stop()
				return nil
			}
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
