package main

import (
	"log"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/keymaps"
	"github.com/phdah/lazydbrix/internal/tui"
    "github.com/phdah/lazydbrix/internal/utils"
)

func main() {
    // Variable decleration
    // TODO: The profile should be set dynamically from the config file
    // profiles := []string{"test", "prod"}
    configPath := "~/.databrickscfg"
    profiles, err := utils.GetProfiles(configPath)
    if err != nil {
        log.Fatalf("Failed to fetch profiles: %v", err)
    }
    var mu sync.Mutex

    // Databricks
    nameToIDMap, _, err := databricks.GetClusterNames(profiles[0])
    if err != nil {
        log.Fatalf("Failed to fetch cluster name: %v", err)
    }

    // TUI components
    app := tview.NewApplication()
    envList := tui.EnvListSetup(profiles)
    prevText := tui.PreTextSetup()
    clusterList := tui.ClusterListSetup(&mu, profiles[0], app, nameToIDMap, prevText)

    // Flex components
    leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(envList, 0, 1, true).
        AddItem(clusterList, 0, 1, false)

    // Create a right Flex
    rightFlex := tview.NewFlex().
        AddItem(prevText, 0, 1, true)

    mainFlex := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(leftFlex, 0, 1, true).
        AddItem(rightFlex, 0, 1, false)

    // Frame components
    frame := tview.NewFrame(mainFlex).
        AddText("lazydbrix", true, tview.AlignCenter, tcell.ColorGreen).
        AddText("Lazily deal with Databricks", true, tview.AlignCenter, tcell.ColorWhite).
        AddText("www.github.com/phdah/lazydbrix", false, tview.AlignRight, tcell.ColorGreen)

    // Set the keymaps
    keymaps.SetKeymaps(app, mainFlex, leftFlex, rightFlex, envList, clusterList, prevText)

    // Set the root and run the application
    if err := app.SetRoot(frame, true).SetFocus(envList).Run(); err != nil {
        panic(err)
    }
}
