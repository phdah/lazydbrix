package utils

import (
	"fmt"
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
	Index       int
	Profile     string
	ClusterName string
	ClusterID   string
}

// Get current cluster from list
func NewClusterFromList(envList *tview.List, clusterList *tview.List) *ClusterFromList {
	clusterIndex := clusterList.GetCurrentItem()
	envMainText, _ := envList.GetItemText(envList.GetCurrentItem())
	clusterMainText, clusterSecondaryText := clusterList.GetItemText(clusterIndex)

	return &ClusterFromList{clusterIndex, envMainText, clusterMainText, clusterSecondaryText}
}

// Helper function to make selections in a list
func ListSelection(envList *tview.List, clusterList *tview.List) *ClusterFromList {
	clusterFromList := NewClusterFromList(envList, clusterList)

	coloredItemFirstText := fmt.Sprintf("[green]%s", clusterFromList.ClusterName)
	clusterList.SetItemText(clusterFromList.Index, coloredItemFirstText, clusterFromList.ClusterID)
	return clusterFromList
}
