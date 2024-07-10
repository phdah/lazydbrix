package main

import (
    "fmt"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "github.com/phdah/lazydbrix/internal/keymaps"
)

func main() {
    app := tview.NewApplication()

    // Create right side lists
    envList := tview.NewList().
        AddItem("dev", "", 0, nil).
        AddItem("prod", "", 0, nil)
    envList.SetBorder(true).SetTitle("Workspaces")

    clusterList := tview.NewList().
        AddItem("Cluster 1", "", 0, nil).
        AddItem("Cluster 2", "", 0, nil).
        AddItem("Cluster 3", "", 0, nil)
    clusterList.SetBorder(true).SetTitle("Clusters")

    // Create a left Flex
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
    fmt.Println(mainFlex.GetItemCount())
    fmt.Println(mainFlex.GetItem(0) == app.GetFocus())
}
