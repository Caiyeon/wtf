package azuredevopsbuilds

import (
	"strings"

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

	azureDevopsBuilds *AzureDevopsBuilds

	Idx int
}

func NewWidget(app *tview.Application) *Widget {
	widget := Widget{
		TextWidget: wtf.NewTextWidget(app, "Azure Devops Builds", "azuredevopsbuilds", false),
		Idx:        0,
	}

	buildIdsRaw := wtf.Config.UString("wtf.mods.azuredevopsbuilds.buildDefinitionIds")
	buildIds := strings.Split(buildIdsRaw, ",")
	for i, buildId := range buildIds {
		buildIds[i] = strings.TrimSpace(buildId)
	}

	widget.azureDevopsBuilds = &AzureDevopsBuilds{
		client:             constructClientFromConfig(),
		buildDefinitionIds: buildIds,
		numBuildsToFetch:   wtf.Config.UInt("wtf.mods.azuredevopsbuilds.numBuildsToFetchPerBuild", 5),
	}

	widget.View.SetInputCapture(widget.keyboardIntercept)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	widget.azureDevopsBuilds.Refresh()
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
