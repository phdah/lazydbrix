package tui

import (
	"sync"

	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/utils"
	"github.com/rivo/tview"
)

type EnvList struct {
	mu          *sync.Mutex
	profile     *string
	app         *tview.Application
	profiles    []string
	clusterList *ClusterList
	dc          *databricks.DatabricksConnection
	detailText  *Text
	List        *tview.List
}

func NewEnvList(mu *sync.Mutex, profile *string, app *tview.Application, profiles []string, clusterList *ClusterList, dc *databricks.DatabricksConnection, detailText *Text) *EnvList {
	return &EnvList{mu, profile, app, profiles, clusterList, dc, detailText, tview.NewList()}
}

func (el *EnvList) Setup() {
	for _, profile := range el.profiles {
		el.List.AddItem(profile, "", 0, nil)
	}

	// Set a function to update the cluster list when the highlighted profile changes
	el.List.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		mainTextUncolored := utils.StripColor(mainText)
		*el.profile = mainTextUncolored
		el.mu.Lock()
		el.mu.Unlock()
		el.clusterList.UpdateClusterList()
	})
	el.List.SetBorder(true).SetTitle("Workspaces")
}
