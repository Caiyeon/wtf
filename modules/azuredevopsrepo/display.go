package azuredevopsrepo

import (
	"fmt"

	"github.com/rivo/tview"
)

func (widget *Widget) display() {
	repo := widget.Repo
	if repo == nil {
		widget.View.SetText(" Azure devops repo data is unavailable")
		return
	}
	if repo.err != nil {
		widget.View.SetText(" Azure devops repo data is unavailable: " + repo.err.Error())
		return
	}

	widget.View.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s - %s", widget.Name, repo.Repo.Name)))

	str := ""
	str = str + " [red]My Pull Requests[white]\n"
	str = str + widget.displayMyCreatedPullRequests(*repo)
	str = str + "\n"
	str = str + " [red]Review Requested[white]\n"
	str = str + widget.displayMyReviewedPullRequests(*repo)
	str = str + "\n"
	str = str + " [red]Open Pull Requests[white]\n"
	str = str + widget.displayOpenPullRequests(*repo)
	str = str + "\n"

	widget.View.SetText(str)
	return
}

func (widget *Widget) displayMyCreatedPullRequests(repo AzureDevopsRepo) string {
	prs := repo.PullRequests

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		if containsUser(widget.User, pr.CreatedBy) {
			str = str + fmt.Sprintf(" [green]%4d[white] %s\n", pr.ID, tview.Escape(pr.Title))
		}
	}

	return str
}

func (widget *Widget) displayMyReviewedPullRequests(repo AzureDevopsRepo) string {
	prs := repo.PullRequests

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		if containsUser(widget.User, pr.Reviewers...) {
			str = str + fmt.Sprintf(" [green]%4d[white] %s\n", pr.ID, tview.Escape(pr.Title))
		}
	}

	return str
}

func (widget *Widget) displayOpenPullRequests(repo AzureDevopsRepo) string {
	prs := repo.PullRequests

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		str = str + fmt.Sprintf(" [green]%4d[white] %s\n", pr.ID, tview.Escape(pr.Title))
	}

	return str
}
