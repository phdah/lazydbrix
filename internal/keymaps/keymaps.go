package keymaps

import (
	"log"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/iface"
	"github.com/phdah/lazydbrix/internal/tui"
	"github.com/phdah/lazydbrix/internal/utils"
	"github.com/rivo/tview"
)

// Set keymaps for a tview.List
func SetEnvKeymaps(app *tview.Application, envList *tview.List) {
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
				app.Stop()
				return nil
			}
		}
		return event
	})
}

// Set keymaps for a tview.List
func SetClusterKeymaps(mu *sync.Mutex, app *tview.Application, envList *tview.List, clusterList *tview.List, detailText *tui.Text, s iface.Selector, dc *databricks.DatabricksConnection) {
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
				app.Stop()
				return nil
			case 's':
				cfl := utils.NewClusterFromList(envList, clusterList)
				dc.ToggleCluster(cfl)
				details, err := dc.GetClusterDetails(cfl)
				if err != nil {
					log.Printf("Failed to get cluster details: %v", err)
				}
				go func() {
					mu.Lock()
					defer mu.Unlock()
					details.UpdateDetails(app, detailText.Text)
				}()
				return nil
			}
		case tcell.KeyEnter:
			cfl := utils.ListSelection(envList, clusterList)
			s.SetSelection(*cfl.GetProfile(), *cfl.GetClusterID(), *cfl.GetClusterName())
			log.Printf("ClusterSelection has been updated to: %s", *s.GetClusterName())
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
