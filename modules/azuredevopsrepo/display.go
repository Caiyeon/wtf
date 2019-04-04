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

	widget.View.SetWrap(false)
	widget.View.SetWordWrap(false)

	if repo.err != nil {
		widget.View.SetWrap(true)
		widget.View.SetWordWrap(true)
		widget.View.SetText(" Azure devops repo data is unavailable: " + repo.err.Error())

		// reset err so that refresh can work again
		repo.err = nil
		return
	}

	widget.View.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s - %s", widget.Name, repo.Repo.Name)))

	prCounter := 0

	str := ""
	str = str + " [red]Created by me[white]\n"
	str = str + widget.displayMyCreatedPullRequests(*repo, &prCounter)
	str = str + "\n"
	str = str + " [red]Assigned to me[white]\n"
	str = str + widget.displayMyReviewedPullRequests(*repo, &prCounter)
	str = str + "\n"
	str = str + " [red]Other Open Pull Requests[white]\n"
	str = str + widget.displayOpenPullRequests(*repo, &prCounter)
	str = str + "\n"

	widget.View.SetText(str)
	return
}

func (widget *Widget) displayMyCreatedPullRequests(repo AzureDevopsRepo, prCounter *int) string {
	prs := repo.PullRequests

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		if *prCounter == widget.maxDisplayedPRs {
			return str
		}

		if containsUser(widget.User, pr.CreatedBy) {
			if *prCounter == widget.SelectedIndex {
				widget.SelectedPR = pr
				str = str + fmt.Sprintf(" [green]%4d[white] %s [yellow]%s[white]\n", pr.ID, prVoteString(pr), tview.Escape(pr.Title))
			} else {
				str = str + fmt.Sprintf(" [green]%4d[white] %s %s\n", pr.ID, prVoteString(pr), tview.Escape(pr.Title))
			}
			*prCounter++
		}
	}

	return str
}

func (widget *Widget) displayMyReviewedPullRequests(repo AzureDevopsRepo, prCounter *int) string {
	prs := repo.PullRequests

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		if *prCounter == widget.maxDisplayedPRs {
			return str
		}

		if containsUser(widget.User, pr.Reviewers...) {
			timeClr, timeString := prTimeString(pr)

			if *prCounter == widget.SelectedIndex {
				widget.SelectedPR = pr
				str = str + fmt.Sprintf(" [green]%4d[white] [lightsalmon]%-8s[white] %s%7s[white] [yellow]%s[yellow]\n",
					pr.ID,
					strings.Split(pr.CreatedBy.DisplayName, " ")[0],
					timeClr,
					timeString,
					tview.Escape(pr.Title),
				)
			} else {
				str = str + fmt.Sprintf(" [green]%4d[white] [lightsalmon]%-8s[white] %s%7s[white] %s\n",
					pr.ID,
					strings.Split(pr.CreatedBy.DisplayName, " ")[0],
					timeClr,
					timeString,
					tview.Escape(pr.Title),
				)
			}
			*prCounter++
		}
	}

	return str
}

func (widget *Widget) displayOpenPullRequests(repo AzureDevopsRepo, prCounter *int) string {
	prs := repo.PullRequests

	if len(prs) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""
	for _, pr := range prs {
		if *prCounter == widget.maxDisplayedPRs {
			return str
		}

		timeClr, timeString := prTimeString(pr)

		if *prCounter == widget.SelectedIndex {
			widget.SelectedPR = pr
			str = str + fmt.Sprintf(" [green]%4d[white] [lightsalmon]%-8s[white] %s%7s[white] [yellow]%s[white]\n",
				pr.ID,
				strings.Split(pr.CreatedBy.DisplayName, " ")[0],
				timeClr,
				timeString,
				tview.Escape(pr.Title),
			)
		} else {
			str = str + fmt.Sprintf(" [green]%4d[white] [lightsalmon]%-8s[white] %s%7s[white] %s\n",
				pr.ID,
				strings.Split(pr.CreatedBy.DisplayName, " ")[0],
				timeClr,
				timeString,
				tview.Escape(pr.Title),
			)
		}
		*prCounter++
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

func prVoteString(pr PullRequest) string {
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
