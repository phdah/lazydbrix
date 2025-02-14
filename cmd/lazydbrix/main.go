package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/phdah/lazydbrix/internal/colors"
	"github.com/phdah/lazydbrix/internal/databricks"
	"github.com/phdah/lazydbrix/internal/keymaps"
	"github.com/phdah/lazydbrix/internal/tui"
	"github.com/phdah/lazydbrix/internal/utils"
)

func main() {
	// Input flags
	debug := flag.Bool("debug", false, "(bool) Flag to run in debug. Default as false")
	outputPath := flag.String("nvim", "", "(string) Path to file to which the cluster selection is written")

	flag.Parse()

	// Discard log output unless -debug is passed
	log.SetOutput(io.Discard)

	// Variable declaration
	configPath := "~/.databrickscfg"
	profiles := utils.GetProfiles(configPath)
	currentProfile := profiles[0]
	var mu sync.Mutex

	// Databricks
	dc := databricks.NewDatabricksConnection(profiles)
	dc.SetWorkspaces()
	dc.SetClusters()

	// TUI components
	app := tview.NewApplication()
	// Set the background color to default
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// Construct tview objects
	detailText := tui.NewText()
	detailText.Setup("Cluster information")
	clusterList := tui.NewClusterList(&mu, &currentProfile, app, dc, detailText)
	clusterList.Setup()
	envList := tui.NewEnvList(&mu, &currentProfile, app, profiles, clusterList, dc, detailText)
	envList.Setup()

	// Set list custom colors
	colors.SetCustomListColor(envList.List)
	colors.SetCustomListColor(clusterList.List)

	// Create a left Flex
	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(envList.List, 0, 1, true).
		AddItem(clusterList.List, 0, 1, false)

	// Create a right Flex
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.AddItem(detailText.Text, 0, 1, true)

	// Pipe logs to tui if debug flag
	if *debug == true {
		logs := tui.NewText()
		logs.Setup("Logs")

		// Redirect standard log output to the TextView
		log.SetOutput(logs.Text)
		rightFlex.AddItem(logs.Text, 0, 3, false)
	}

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFlex, 0, 1, true).
		AddItem(rightFlex, 0, 1, false)

	// Frame components
	frame := tview.NewFrame(mainFlex).
		AddText("lazydbrix", true, tview.AlignCenter, tcell.ColorGreen).
		AddText("Lazily deal with Databricks clusters", true, tview.AlignCenter, tcell.ColorWhite).
		AddText("www.github.com/phdah/lazydbrix", false, tview.AlignRight, tcell.ColorGreen).
		AddText("Quit: q | Select cluster: <enter> | Toggle cluster: s", false, tview.AlignLeft, tcell.ColorBlue)

	// Setup a cluster selection struct
	cs := tui.NewClusterSelection()

	// Set the keymaps
	keymaps.SetEnvKeymaps(app, envList.List)
	keymaps.SetClusterKeymaps(&mu, app, envList.List, clusterList.List, detailText, cs, dc)
	keymaps.SetFlexKeymaps(app, leftFlex)
	keymaps.SetFlexKeymaps(app, rightFlex)
	keymaps.SetMainFlexKeymaps(app, mainFlex)

	// Set the root and run the application
	if err := app.SetRoot(frame, true).SetFocus(envList.List).Run(); err != nil {
		panic(err)
	}

	if cs.ClusterID != "" {
		if *outputPath != "" {
			clusterSelectedLua := fmt.Sprintf("let $DATABRICKS_CONFIG_PROFILE =\"%s\"\nlet $CLUSTER_NAME = \"%s\"\nlet $DATABRICKS_CLUSTER_ID = \"%s\"\n", *cs.GetProfile(), *cs.GetClusterName(), *cs.GetClusterID())

			file, createErr := os.Create(*outputPath)
			if createErr != nil {
				fmt.Println("Error creating to file:", createErr)
				return
			}
			_, writingErr := file.WriteString(clusterSelectedLua)
			if writingErr != nil {
				fmt.Println("Error writing to file:", writingErr)
			}
			fmt.Println("Successfully updated cluster info in file:", *outputPath)
		} else {
			fmt.Printf("Cluster selected, but not set: {Profile: %s, ClusterName: %s, ClusterID: %s}\n", *cs.GetProfile(), *cs.GetClusterName(), *cs.GetClusterID())
		}
	} else {
		fmt.Println("No selection made")
	}
}
