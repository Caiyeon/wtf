package azuredevopsrepo

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/wtf"
)

const HelpText = `
  Keyboard commands for GitHub:

    /: Show/hide this help window
	r: Refresh the data

    arrow up:   Previous pull request
    arrow down: Next pull request

    return: Open the selected pull request in a browser
`

type Widget struct {
	wtf.HelpfulWidget
	wtf.TextWidget

	Repo          *AzureDevopsRepo
	User          string
	SelectedIndex int
	SelectedPR    PullRequest

	maxDisplayedPRs int
}

func NewWidget(app *tview.Application, pages *tview.Pages) *Widget {
	widget := Widget{
		HelpfulWidget: wtf.NewHelpfulWidget(app, pages, HelpText),
		TextWidget:    wtf.NewTextWidget(app, "Azure Devops Repo", "azuredevopsrepo", true),
		SelectedIndex: -1,
	}

	widget.User = wtf.Config.UString("wtf.mods.azuredevopsrepo.user")
	widget.maxDisplayedPRs = wtf.Config.UInt("wtf.mods.azuredevopsrepo.pullRequestCount", 10)
	widget.Repo = NewRepo(wtf.Config.UString("wtf.mods.azuredevopsrepo.repository"))

	widget.HelpfulWidget.SetView(widget.View)
	widget.View.SetDoneFunc(widget.doneFunc)
	widget.View.SetInputCapture(widget.keyboardIntercept)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	// if user isn't focused, don't highlight any prs
	if !widget.View.HasFocus() {
		widget.SelectedIndex = -1
		widget.SelectedPR = PullRequest{}
	}
	widget.Repo.Refresh()
	widget.display()
}

func (widget *Widget) Prev() {
	widget.SelectedIndex = widget.SelectedIndex - 1
	if widget.SelectedIndex < 0 {
		widget.SelectedIndex = widget.maxDisplayedPRs - 1
	}
	widget.display()
}

func (widget *Widget) Next() {
	widget.SelectedIndex = widget.SelectedIndex + 1
	if widget.SelectedIndex == widget.maxDisplayedPRs {
		widget.SelectedIndex = 0
	}
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) openSelectedPR() {
	URL := widget.Repo.Repo.RemoteUrl + "/pullrequest"
	if widget.SelectedIndex != -1 {
		URL = URL + "/" + strconv.Itoa(widget.SelectedPR.ID)
	}
	wtf.OpenFile(URL)
}

func (widget *Widget) doneFunc(event tcell.Key) {
	widget.SelectedIndex = -1
	widget.SelectedPR = PullRequest{}
	widget.display()
}

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
	case tcell.KeyEnter:
		widget.openSelectedPR()
		return nil
	case tcell.KeyUp:
		widget.Prev()
		return nil
	case tcell.KeyDown:
		widget.Next()
		return nil
	case tcell.KeyEscape:
		widget.SelectedIndex = 0
		return nil
	default:
		return event
	}
}
