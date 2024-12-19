package tui

import (
	"sync"

	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/utils"
	"github.com/rivo/tview"
)

func EnvListSetup(mu *sync.Mutex, profile *string, app *tview.Application, profiles []string, clusterList *tview.List, dc *databricks.DatabricksConnection, prevText *tview.TextView) *tview.List {
	envList := tview.NewList()
	for _, profile := range profiles {
		envList.AddItem(profile, "", 0, nil)
	}

	// Set a function to update the cluster list when the highlighted profile changes
	envList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		mainTextUncolored := utils.StripColor(mainText)
		*profile = mainTextUncolored
		mu.Lock()
		nameToIDMap := dc.ProfileClusters[*profile]
		mu.Unlock()
		UpdateClusterList(mu, app, profile, clusterList, nameToIDMap, prevText)
	})
	envList.SetBorder(true).SetTitle("Workspaces")

	return envList
}
