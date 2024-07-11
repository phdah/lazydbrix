package main

import (
    "log"
    "fmt"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"

    "github.com/phdah/lazydbrix/internal/databricks"
    "github.com/phdah/lazydbrix/internal/keymaps"
)

func main() {
    // Variable decleration
    // TODO: The profile should be set dynamically
    profile := "test"

    // Databricks
    nameToIDMap, _, err := databricks.GetClusterNames(profile)
    if err != nil {
        log.Fatalf("Failed to fetch cluster name: %v", err)
    }

    // TUI
    app := tview.NewApplication()

    // Create right side lists
    envList := tview.NewList().
        AddItem("dev", "", 0, nil).
        AddItem("prod", "", 0, nil)

    clusterList := tview.NewList()

    for clusterName, clusterId := range nameToIDMap {
        clusterList.AddItem(fmt.Sprintf("%s: %s", clusterName, clusterId), "", 0, nil)
    }

    // Create a left Flex
    envList.SetBorder(true).SetTitle("Workspaces")
    clusterList.SetBorder(true).SetTitle("Clusters")
    leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(envList, 0, 1, true).
        AddItem(clusterList, 0, 1, false)

    // Create a right Flex
    prevBox := tview.NewBox().
        SetBorder(true).
        SetTitle("Cluster information")

    rightFlex := tview.NewFlex().
        AddItem(prevBox, 0, 1, true)

    mainFlex := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(leftFlex, 0, 1, true).
        AddItem(rightFlex, 0, 1, false)

    // Create the frame and add text to it
    frame := tview.NewFrame(mainFlex).
        AddText("lazydbrix", true, tview.AlignCenter, tcell.ColorGreen).
        AddText("Lazily deal with Databricks", true, tview.AlignCenter, tcell.ColorWhite).
        AddText("www.github.com/phdah/lazydbrix", false, tview.AlignRight, tcell.ColorGreen)

    // Set the keymaps
    keymaps.SetKeymaps(app, mainFlex, leftFlex, rightFlex, envList, clusterList, prevBox)

    // Set the root and run the application
    if err := app.SetRoot(frame, true).SetFocus(envList).Run(); err != nil {
        panic(err)
    }
}
