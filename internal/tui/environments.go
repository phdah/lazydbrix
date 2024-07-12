package tui

import (
    "sync"
    "log"

    "github.com/rivo/tview"
    "github.com/elliotchance/orderedmap/v2"
)

func EnvListSetup(mu *sync.Mutex, profile *string, app *tview.Application, profiles []string, clusterList *tview.List, allNameToIDMap map[string]*orderedmap.OrderedMap[string, string], prevText *tview.TextView) (*tview.List) {
    envList := tview.NewList()
    for _, profile := range profiles {
        envList.AddItem(profile, "", 0, nil)
    }

    // Set a function to update the cluster list when the highlighted profile changes
    envList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
        log.Printf("Updateing %s to %s", *profile, mainText)
        *profile = mainText
        mu.Lock()
        nameToIDMap := allNameToIDMap[*profile]
        mu.Unlock()
        log.Printf("New nameToIdMap[%s]", *profile)
        UpdateClusterList(mu, app, profile, clusterList, nameToIDMap, prevText)
    })
    envList.SetBorder(true).SetTitle("Workspaces")

    return envList
}
