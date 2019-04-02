package azuredevopsrepo

import (
	"fmt"
	"strings"
	"time"

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
	str = str + " [red]Created by me[white]\n"
	str = str + widget.displayMyCreatedPullRequests(*repo)
	str = str + "\n"
	str = str + " [red]Assigned to me[white]\n"
	str = str + widget.displayMyReviewedPullRequests(*repo)
	str = str + "\n"
	str = str + " [red]Other Open Pull Requests[white]\n"
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
			str = str + fmt.Sprintf(" [green]%4d[white] %s %s\n", pr.ID, reviewString(pr), tview.Escape(pr.Title))
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
			timeClr, timeString := prTimeString(pr)
			str = str + fmt.Sprintf(" [green]%4d[white] [lightsalmon]%-8s[white] %s%7s[white] %s\n",
				pr.ID,
				strings.Split(pr.CreatedBy.DisplayName, " ")[0],
				timeClr,
				timeString,
				tview.Escape(pr.Title),
			)
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
		timeClr, timeString := prTimeString(pr)
		str = str + fmt.Sprintf(" [green]%4d[white] [lightsalmon]%-8s[white] %s%7s[white] %s\n",
			pr.ID,
			strings.Split(pr.CreatedBy.DisplayName, " ")[0],
			timeClr,
			timeString,
			tview.Escape(pr.Title),
		)
	}

	return str
}

func prTimeString(pr PullRequest) (string, string) {
	t, err := time.Parse(time.RFC3339Nano, pr.Created)
	if err != nil {
		return "", ""
	}

	if time.Since(t) < 1*time.Hour {
		return "[orange]", fmt.Sprintf("%.0fm ago", time.Since(t).Minutes())
	} else if time.Since(t) < 24*time.Hour {
		return "[orange]", fmt.Sprintf("%.0fh ago", time.Since(t).Hours())
	} else {
		return "[grey]", t.Format("Jan 02")
	}
}

func reviewString(pr PullRequest) string {
	s := ""
	for _, reviewer := range pr.Reviewers {
		if !reviewer.IsContainer {
			if reviewer.Vote == 10 || reviewer.Vote == 5 {
				s = s + "[green]✓[white]"
			}
			if reviewer.Vote == -10 || reviewer.Vote == -5 {
				s = s + "[orange]✗[white]"
			}
		}
	}
	s = "|" + s + "|"
	return s
}
