package tui

import (
    "sync"
    "log"

	"github.com/rivo/tview"

	"github.com/phdah/lazydbrix/internal/databricks"
)

func ClusterListSetup(mu *sync.Mutex, profile string, app *tview.Application, nameToIDMap map[string]string, prevText *tview.TextView) (*tview.List) {
    clusterList := tview.NewList()

    var firstClusterID string
    for clusterName, clusterID := range nameToIDMap {
        if firstClusterID == "" {
            firstClusterID = clusterID
        }
        clusterList.AddItem(clusterName, "", 0, nil)
    }

	// Set a function to update the preview text view when the highlighted item changes
	clusterList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		go func() {
            clusterID := nameToIDMap[mainText]
            details, err := databricks.GetClusterDetails(profile, clusterID)
            if err != nil {
                log.Fatalf("Failed to fetch cluster details: %v", err)
            }
            // mutex ensures thread safety for the goroutine
            mu.Lock()
            defer mu.Unlock()
            app.QueueUpdateDraw(func() {
                prevText.SetText(databricks.FormatClusterDetails(details))
            })
        }()
	})

    clusterList.SetBorder(true).SetTitle("Clusters")

    details, _ := databricks.GetClusterDetails(profile, firstClusterID)
    prevText.SetText(databricks.FormatClusterDetails(details))

    return clusterList
}
