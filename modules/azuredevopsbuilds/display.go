package azuredevopsbuilds

import (
	"fmt"
	"strings"
	"time"
)

func (widget *Widget) display() {
	azureDevopsBuilds := widget.azureDevopsBuilds

	if azureDevopsBuilds == nil {
		widget.View.SetText(" Azure devops data is unavailable")
		return
	}
	if azureDevopsBuilds.err != nil {
		widget.View.SetText(" " + azureDevopsBuilds.err.Error())
		return
	}

	widget.View.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s", widget.Name)))

	str := ""
	str = str + displayBuilds(azureDevopsBuilds.buildDefinitionIds, azureDevopsBuilds.builds)
	str = str + "\n"

	widget.View.SetText(str)
	return
}

func displayBuilds(buildIds []string, allBuilds map[string][]Build) string {
	if len(buildIds) == 0 || len(allBuilds) == 0 {
		return " [grey]none[white]\n"
	}

	str := ""

	// need to loop over the slice to preserve ordering of sets of builds
	for _, buildId := range buildIds {
		if len(allBuilds[buildId]) == 0 {
			continue
		}

		str = str + " [red]" + allBuilds[buildId][0].Definition.Name + "[white]\n"

		for _, build := range allBuilds[buildId] {
			progClr, progString := progressString(build)
			timeClr, timeString := timeString(build)

			str = str + fmt.Sprintf(" [green]%7d[white] %s%-11s[white] [lightsalmon]%-8s[white] %s%-7s[white] %s\n",
				build.ID,
				progClr,
				progString,
				strings.Split(build.RequestedFor.DisplayName, " ")[0],
				timeClr,
				timeString,
				build.BuildNumber,
			)
		}

		str = str + "\n"
	}

	return str
}

func timeString(b Build) (string, string) {
	tRaw := ""
	tZone := time.FixedZone("UTC-8", -8*60*60)

	switch b.Status {
	case "inProgress":
		tRaw = b.StartTime
		t, err := time.Parse(time.RFC3339Nano, tRaw)
		if err != nil {
			return "", ""
		}
		return "[orange]", fmt.Sprintf("%.0fm ago", time.Since(t).Minutes())

	case "completed":
		tRaw = b.FinishTime
		t, err := time.ParseInLocation(time.RFC3339Nano, tRaw, tZone)
		if err != nil {
			return "", ""
		}
		return "[grey]", t.Format(time.Kitchen)

	default:
		return "Pending", ""
	}
}

func progressString(b Build) (string, string) {
	if b.Status == "inProgress" {
		return "[yellow]", "In Progress"
	}
	if b.Status == "notStarted" {
		return "[grey]", "Not Started"
	}

	switch b.Result {
	case "succeeded":
		return "[lightgreen]", "Succeeded"

	case "partiallySucceeded":
		fallthrough
	case "failed":
		return "[red]", "Failed"

	case "canceled":
		return "[grey]", "Canceled"

	case "none":
		return "[grey]", "None"
	}

	return "", ""
}
