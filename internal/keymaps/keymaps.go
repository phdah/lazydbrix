package keymaps

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "github.com/phdah/lazydbrix/internal/utils"
)

func SetKeymaps(app *tview.Application, mainFlex *tview.Flex, leftFlex, rightFlex *tview.Flex, envList, clusterList *tview.List, prevBox *tview.Box) {
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyTab: // Handle Tab key
            utils.MoveFlexRight(app, mainFlex)
            return nil
        case tcell.KeyBacktab: // Handle Shift+Tab key
            utils.MoveFlexLeft(app, mainFlex)
            return nil
        case tcell.KeyRune: // Handle 'j', 'k', 'h', 'l' using Rune
            switch event.Rune() {
            case 'j': // Move down in the list
                if app.GetFocus() == envList {
                    utils.MoveListDown(envList)
                } else if app.GetFocus() == clusterList {
                    utils.MoveListDown(clusterList)
                }
                return nil
            case 'k': // Move up in the list
                if app.GetFocus() == envList {
                    utils.MoveListUp(envList)
                } else if app.GetFocus() == clusterList {
                    utils.MoveListUp(clusterList)
                }
                return nil
            case 'h': // Move focus to the list underneath in the flex
                if app.GetFocus() == rightFlex || app.GetFocus() == prevBox {
                    utils.MoveFlexItemUp(app, rightFlex)
                } else if app.GetFocus() == leftFlex || app.GetFocus() == envList || app.GetFocus() == clusterList {
                    utils.MoveFlexItemUp(app, leftFlex)
                }
                return nil
            case 'l': // Move focus to the list above in the flex
                if app.GetFocus() == rightFlex || app.GetFocus() == prevBox {
                    utils.MoveFlexItemDown(app, rightFlex)
                } else if app.GetFocus() == leftFlex || app.GetFocus() == envList || app.GetFocus() == clusterList {
                    utils.MoveFlexItemDown(app, leftFlex)
                }
                return nil
            }
        }
        return event // Forward to other handlers
    })
}
