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

	str := ""
	if wtf.Config.UBool("wtf.mods.azuredevopsworkitems.displayName") {
		for _, item := range items {
			if item.Fields.State == "In Progress" {
				str = str + fmt.Sprintf(" [yellow]%s[white] [green]%7d[white] [orange]%s[white] %s\n",
					item.Fields.WorkItemType, item.ID, item.Fields.AssignedTo.DisplayName, tview.Escape(item.Fields.Title),
				)
			}
		}
	} else {
		for _, item := range items {
			if item.Fields.State == "In Progress" {
				str = str + fmt.Sprintf(" [yellow]%s[white] [green]%7d[white] %s\n",
					item.Fields.WorkItemType, item.ID, tview.Escape(item.Fields.Title),
				)
			}
		}
	}

	return str
}

func displayToDoWorkItems(items []WorkItem) string {
	if len(items) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	if wtf.Config.UBool("wtf.mods.azuredevopsworkitems.displayName") {
		for _, item := range items {
			if item.Fields.State == "To Do" {
				str = str + fmt.Sprintf(" [yellow]%s[white] [green]%7d[white] [orange]%s[white] %s\n",
					item.Fields.WorkItemType, item.ID, item.Fields.AssignedTo.DisplayName, tview.Escape(item.Fields.Title),
				)
			}
		}
	} else {
		for _, item := range items {
			if item.Fields.State == "To Do" {
				str = str + fmt.Sprintf(" [yellow]%s[white] [green]%7d[white] %s\n",
					item.Fields.WorkItemType, item.ID, tview.Escape(item.Fields.Title),
				)
			}
		}
	}

	return str
}

func displayOtherWorkItems(items []WorkItem) string {
	if len(items) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	if wtf.Config.UBool("wtf.mods.azuredevopsworkitems.displayName") {
		for _, item := range items {
			if item.Fields.State != "In Progress" && item.Fields.State != "To Do" {
				str = str + fmt.Sprintf(" [orange]%s[white] [yellow]%s[white] [green]%7d[white] [lightsalmon]%s[white] %s\n",
					item.Fields.State, item.Fields.WorkItemType, item.ID, item.Fields.AssignedTo.DisplayName, tview.Escape(item.Fields.Title),
				)
			}
		}
	} else {
		for _, item := range items {
			if item.Fields.State != "In Progress" && item.Fields.State != "To Do" {
				str = str + fmt.Sprintf(" [orange]%s[white] [yellow]%s[white] [green]%7d[white] %s\n",
					item.Fields.State, item.Fields.WorkItemType, item.ID, tview.Escape(item.Fields.Title),
				)
			}
		}
	}

	return str
}
