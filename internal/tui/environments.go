package tui

import (
	"github.com/rivo/tview"
)

func EnvListSetup(profiles []string) (*tview.List) {
    envList := tview.NewList()
    for _, profile := range profiles {
        envList.AddItem(profile, "", 0, nil)
    }

    envList.SetBorder(true).SetTitle("Workspaces")

    return envList
}
