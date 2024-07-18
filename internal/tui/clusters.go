package tui

import (
	"log"
	"sync"

	"github.com/rivo/tview"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/utils"
)

func ClusterListSetup(mu *sync.Mutex, profile *string, app *tview.Application, allNameToIDMap map[string]*orderedmap.OrderedMap[string, string], prevText *tview.TextView) *tview.List {
	clusterList := tview.NewList()
	clusterList.ShowSecondaryText(false)
	nameToIDMap := allNameToIDMap[*profile]

	var firstClusterID string
	for c := nameToIDMap.Front(); c != nil; c = c.Next() {
		if firstClusterID == "" {
			firstClusterID = c.Value
		}
		clusterList.AddItem(c.Key, c.Value, 0, nil)
	}

	// Set a function to update the preview text view when the highlighted item changes
	clusterList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
        mainTextUncolored := utils.StripColor(mainText)
		log.Printf("->clusterList: profile %s, cluster %s", *profile, mainTextUncolored)
		go func() {
			nameToIDMap = allNameToIDMap[*profile]
			clusterID := nameToIDMap.GetElement(mainTextUncolored).Value
			details, err := databricks.GetClusterDetails(profile, clusterID)
			if err != nil {
				log.Printf("Failed to fetch cluster details: %v", err)
			}
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

// UpdateClusterList updates the cluster list based on the selected profile
func UpdateClusterList(mu *sync.Mutex, app *tview.Application, profile *string, clusterList *tview.List, nameToIDMap *orderedmap.OrderedMap[string, string], prevText *tview.TextView) {
	mu.Lock()
	go func() {
		defer mu.Unlock()
		clusterList.Clear()
		var firstClusterID string
		for c := nameToIDMap.Front(); c != nil; c = c.Next() {
			if firstClusterID == "" {
				firstClusterID = c.Value
			}
			clusterList.AddItem(c.Key, c.Value, 0, nil)
		}
		log.Printf("UpdateClusterList, profile: %s. firstClusterID is %s", *profile, firstClusterID)
	}()
}
