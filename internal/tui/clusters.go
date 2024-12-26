package tui

import (
	"log"
	"sync"

	"github.com/rivo/tview"

	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/utils"
)

type ClusterList struct {
	mu         *sync.Mutex
	profile    *string
	app        *tview.Application
	dc         *databricks.DatabricksConnection
	detailText *Text
	List       *tview.List
}

func NewClusterList(mu *sync.Mutex, profile *string, app *tview.Application, dc *databricks.DatabricksConnection, detailText *Text) *ClusterList {
	return &ClusterList{mu, profile, app, dc, detailText, tview.NewList()}
}

func (cl *ClusterList) Setup() {
	cl.List.ShowSecondaryText(false)
	nameToIDMap := cl.dc.ProfileClusters[*cl.profile]

	var firstClusterID string
	for c := nameToIDMap.Front(); c != nil; c = c.Next() {
		if firstClusterID == "" {
			firstClusterID = c.Value
		}
		cl.List.AddItem(c.Key, c.Value, 0, nil)
	}

	// Set a function to update the preview text view when the highlighted item changes
	cl.List.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		mainTextUncolored := utils.StripColor(mainText)
		log.Printf("->clusterList: profile %s, cluster %s", *cl.profile, mainTextUncolored)
		go func() {
			nameToIDMap = cl.dc.ProfileClusters[*cl.profile]
			clusterID := nameToIDMap.GetElement(mainTextUncolored).Value
			details, err := cl.dc.GetClusterDetails(databricks.NewCluster(*cl.profile, clusterID, mainTextUncolored))
			if err != nil {
				log.Printf("Failed to fetch cluster details: %v", err)
			}
			cl.mu.Lock()
			defer cl.mu.Unlock()
			details.UpdateDetails(cl.app, cl.detailText.Text)
			log.Println("Done updateing text field")
		}()
	})

	cl.List.SetBorder(true).SetTitle("Clusters")

	details, _ := cl.dc.GetClusterDetails(databricks.NewCluster(*cl.profile, firstClusterID, ""))
	cl.detailText.Text.SetText(details.Format())
}

// UpdateClusterList updates the cluster list based on the selected profile
func (cl *ClusterList) UpdateClusterList() {
	cl.mu.Lock()
	go func() {
		defer cl.mu.Unlock()
		nameToIDMap := cl.dc.ProfileClusters[*cl.profile]
		cl.List.Clear()
		var firstClusterID string
		for c := nameToIDMap.Front(); c != nil; c = c.Next() {
			if firstClusterID == "" {
				firstClusterID = c.Value
			}
			cl.List.AddItem(c.Key, c.Value, 0, nil)
		}
		log.Printf("UpdateClusterList, profile: %s. firstClusterID is %s", *cl.profile, firstClusterID)
	}()
}
