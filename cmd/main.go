package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/keymaps"
	"github.com/phdah/lazydbrix/internal/tui"
	"github.com/phdah/lazydbrix/internal/utils"
)

func main() {
	// Input flags
	debug := flag.Bool("debug", false, "(bool) Flag to run in debug. Default as false")
	outputPath := flag.String("output", "", "(string) Path to file to which the cluster selection is written")

	flag.Parse()
	// Variable declaration
	configPath := "~/.databrickscfg"
	profiles := utils.GetProfiles(configPath)
	currentProfile := profiles[0]
	var mu sync.Mutex

	logs := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	logs.SetBorder(true).SetTitle("Logs")

	// Redirect standard log output to the TextView
	log.SetOutput(logs)

	// Databricks
	allNameToIDMap := databricks.GetAllEnvClusters(&mu, profiles)

	// TUI components
	app := tview.NewApplication()
	prevText := tui.PreTextSetup()
	clusterList := tui.ClusterListSetup(&mu, &currentProfile, app, allNameToIDMap, prevText)
	envList := tui.EnvListSetup(&mu, &currentProfile, app, profiles, clusterList, allNameToIDMap, prevText)

	// Flex components
	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(envList, 0, 1, true).
		AddItem(clusterList, 0, 1, false)

	// Create a right Flex
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.AddItem(prevText, 0, 1, true)

	if *debug == true {
		rightFlex.AddItem(logs, 0, 3, false)
	}

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
	keymaps.SetListKeymaps(envList)
	keymaps.SetListKeymaps(clusterList)
	keymaps.SetFlexKeymaps(app, leftFlex)
	keymaps.SetFlexKeymaps(app, rightFlex)
	keymaps.SetMainFlexKeymaps(app, mainFlex)

	// Set the root and run the application
	if err := app.SetRoot(frame, true).SetFocus(envList).Run(); err != nil {
		panic(err)
	}

	if *outputPath != "" {
		envMainText, _ := envList.GetItemText(envList.GetCurrentItem())
		clusterMainText, clusterSecondaryText := clusterList.GetItemText(clusterList.GetCurrentItem())
		clusterSelection := fmt.Sprintf("{\"PROFILE\": \"%s\", \"CLUSTER_NAME\": \"%s\", \"CLUSTER_ID\": \"%s\"}", envMainText, clusterMainText, clusterSecondaryText)

		file, createErr := os.Create(*outputPath)
		if createErr != nil {
			fmt.Println("Error creating to file:", createErr)
			return
		}
		_, writingErr := file.WriteString(clusterSelection)
		if writingErr != nil {
			fmt.Println("Error writing to file:", writingErr)
		}
		fmt.Println("Successfully updated cluster info in file:", *outputPath)
	}
}
