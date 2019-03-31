package azuredevopsrepo

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/wtf"
)

const HelpText = `
  Keyboard commands for GitHub:

    /: Show/hide this help window
    r: Refresh the data

    return: Open the selected repository in a browser
`

type Widget struct {
	wtf.HelpfulWidget
	wtf.TextWidget

	Repo *AzureDevopsRepo
	User string
	Idx  int
}

func NewWidget(app *tview.Application) *Widget {
	widget := Widget{
		TextWidget: wtf.NewTextWidget(app, "Azure Devops Repo", "azuredevopsrepo", false),
		Idx:        0,
	}

	widget.User = wtf.Config.UString("wtf.mods.azuredevopsrepo.user")

	widget.Repo = NewRepo(wtf.Config.UString("wtf.mods.azuredevopsrepo.repository"))

	widget.View.SetInputCapture(widget.keyboardIntercept)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.Repo.Refresh()
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) keyboardIntercept(event *tcell.EventKey) *tcell.EventKey {
	switch string(event.Rune()) {
	case "/":
		widget.ShowHelp()
		return nil
	case "r":
		widget.Refresh()
		return nil
	}

	switch event.Key() {
	default:
		return event
	}
}
