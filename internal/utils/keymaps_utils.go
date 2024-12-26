package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/rivo/tview"
)

// Helper functions to move up and down in the list
func MoveListDown(list *tview.List) {
	currentItem := list.GetCurrentItem()
	if currentItem < list.GetItemCount()-1 {
		list.SetCurrentItem(currentItem + 1)
	}
}

func MoveListUp(list *tview.List) {
	currentItem := list.GetCurrentItem()
	if currentItem > 0 {
		list.SetCurrentItem(currentItem - 1)
	}
}

// Helper functions to move focus up and down in the flex
func MoveFlexItemUp(app *tview.Application, flex *tview.Flex) {
	focused := app.GetFocus()
	for i := 0; i < flex.GetItemCount(); i++ {
		item := flex.GetItem(i)
		if item == focused && i > 0 {
			app.SetFocus(flex.GetItem(i - 1))
			break
		}
	}
}

func MoveFlexItemDown(app *tview.Application, flex *tview.Flex) {
	focused := app.GetFocus()
	for i := 0; i < flex.GetItemCount(); i++ {
		item := flex.GetItem(i)
		if item == focused && i < flex.GetItemCount()-1 {
			app.SetFocus(flex.GetItem(i + 1))
			break
		}
	}
}

// Helper functions to move focus left and right in the flex
func MoveFlexRight(app *tview.Application, flex *tview.Flex) {
	for i := 0; i < flex.GetItemCount(); i++ {
		item := flex.GetItem(i)
		if item.HasFocus() && i < flex.GetItemCount()-1 {
			app.SetFocus(flex.GetItem(i + 1))
			break
		}
	}
}

func MoveFlexLeft(app *tview.Application, flex *tview.Flex) {
	for i := 0; i < flex.GetItemCount(); i++ {
		item := flex.GetItem(i)
		if item.HasFocus() && i > 0 {
			app.SetFocus(flex.GetItem(i - 1))
			break
		}
	}
}

type ClusterFromList struct {
	mu          sync.Mutex
	Index       int
	Profile     string
	ClusterName string
	ClusterID   string
}

// Get current cluster from list
func NewClusterFromList(envList *tview.List, clusterList *tview.List) *ClusterFromList {
	clusterIndex := clusterList.GetCurrentItem()
	profile, _ := envList.GetItemText(envList.GetCurrentItem())
	clusterName, clusterID := clusterList.GetItemText(clusterIndex)

	cfl := &ClusterFromList{}
	cfl.SetClusterIndex(clusterIndex)
	cfl.SetSelection(profile, clusterID, clusterName)

	return cfl
}

func (cfl *ClusterFromList) GetProfile() *string {
	cfl.mu.Lock()
	defer cfl.mu.Unlock()
	return &cfl.Profile
}
func (cfl *ClusterFromList) GetClusterID() *string {
	cfl.mu.Lock()
	defer cfl.mu.Unlock()
	return &cfl.ClusterID
}
func (cfl *ClusterFromList) GetClusterName() *string {
	cfl.mu.Lock()
	defer cfl.mu.Unlock()
	return &cfl.ClusterName
}
func (cfl *ClusterFromList) GetIndex() int {
	cfl.mu.Lock()
	defer cfl.mu.Unlock()
	return cfl.Index
}
func (cfl *ClusterFromList) SetClusterIndex(index int) {
	cfl.mu.Lock()
	defer cfl.mu.Unlock()
	cfl.Index = index
}
func (cfl *ClusterFromList) SetSelection(profile string, clusterID string, clusterName string) {
	cfl.mu.Lock()
	defer cfl.mu.Unlock()
	log.Printf("Cluster selected: %s (ID: %s, Profile: %s)", clusterName, clusterID, profile)
	cfl.Profile = profile
	cfl.ClusterID = clusterID
	cfl.ClusterName = clusterName
}

// Helper function to make selections in a list
func ListSelection(envList *tview.List, clusterList *tview.List) *ClusterFromList {
	cfl := NewClusterFromList(envList, clusterList)

	coloredItemFirstText := fmt.Sprintf("[green]%s", *cfl.GetClusterName())
	clusterList.SetItemText(cfl.GetIndex(), coloredItemFirstText, *cfl.GetClusterID())
	return cfl
}
