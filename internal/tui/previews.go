package tui

import (
    "github.com/rivo/tview"
)

func PreTextSetup() (*tview.TextView) {
    prevText := tview.NewTextView().
        SetDynamicColors(true).
        SetRegions(true).
        SetWrap(true)

    prevText.SetBorder(true).SetTitle("Cluster information")

    return prevText
}
