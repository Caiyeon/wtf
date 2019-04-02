package azuredevopsworkitems

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/wtf"
)

func (widget *Widget) display() {
	azureWorkItems := widget.azureWorkItems

	if azureWorkItems == nil {
		widget.View.SetText(" Azure devops data is unavailable")
		return
	}
	if azureWorkItems.err != nil {
		widget.View.SetText(" " + azureWorkItems.err.Error())
		return
	}

	widget.View.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s - %s", widget.Name, azureWorkItems.queryId)))

	str := ""
	str = str + " [red]In Progress[white]\n"
	str = str + displayInProgressWorkItems(azureWorkItems.workItems)
	str = str + "\n"
	str = str + " [red]To Do[white]\n"
	str = str + displayToDoWorkItems(azureWorkItems.workItems)
	str = str + "\n"
	str = str + " [red]Other[white]\n"
	str = str + displayOtherWorkItems(azureWorkItems.workItems)
	str = str + "\n"

	widget.View.SetText(str)
	return
}

func displayInProgressWorkItems(items []WorkItem) string {
	if len(items) == 0 {
		return " [grey]none[white]\n"
	}

	displayName := wtf.Config.UBool("wtf.mods.azuredevopsworkitems.displayName")

	str := ""
	for _, item := range items {
		if item.Fields.State == "In Progress" || item.Fields.State == "Committed" {
			str = str + fmt.Sprintf(" [green]%7d[white] [yellow]%s[white] ",
				item.ID, item.Fields.WorkItemType,
			)
			if displayName {
				str = str + fmt.Sprintf("[lightsalmon]%s[white] ", item.Fields.AssignedTo.DisplayName)
			}
			str = str + tview.Escape(item.Fields.Title) + "\n"
		}
	}

	return str
}

func displayToDoWorkItems(items []WorkItem) string {
	if len(items) == 0 {
		return " [grey]none[white]\n"
	}

	displayName := wtf.Config.UBool("wtf.mods.azuredevopsworkitems.displayName")

	str := ""
	for _, item := range items {
		if item.Fields.State == "To Do" {
			str = str + fmt.Sprintf(" [green]%7d[white] [yellow]%s[white] ",
				item.ID, item.Fields.WorkItemType,
			)
			if displayName {
				str = str + fmt.Sprintf("[lightsalmon]%s[white] ", item.Fields.AssignedTo.DisplayName)
			}
			str = str + tview.Escape(item.Fields.Title) + "\n"
		}
	}

	return str
}

func displayOtherWorkItems(items []WorkItem) string {
	if len(items) == 0 {
		return " [grey]none[white]\n"
	}

	displayName := wtf.Config.UBool("wtf.mods.azuredevopsworkitems.displayName")

	str := ""
	for _, item := range items {
		if itemIsOther(item) {
			str = str + fmt.Sprintf(" [green]%7d[white] [yellow]%s[white] [orange]%s[white] ",
				item.ID, item.Fields.WorkItemType, item.Fields.State,
			)
			if displayName {
				str = str + fmt.Sprintf("[lightsalmon]%s[white] ", item.Fields.AssignedTo.DisplayName)
			}
			str = str + tview.Escape(item.Fields.Title) + "\n"
		}
	}

	return str
}

func itemIsOther(item WorkItem) bool {
	notOther := []string{
		"In Progress",
		"To Do",
		"Done",
		"Removed",
		"Committed",
	}

	for _, s := range notOther {
		if item.Fields.State == s {
			return false
		}
	}

	return true
}
