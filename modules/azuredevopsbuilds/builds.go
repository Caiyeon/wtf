package azuredevopsbuilds

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	az "github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/wtfutil/wtf/wtf"
)

type AzureDevopsBuilds struct {
	client             *az.Client
	buildDefinitionIds []string
	builds             map[string][]Build
	numBuildsToFetch   int

	err error
}

type BuildsListResponse struct {
	Builds []Build `json:"value"`
}

type Build struct {
	Definition    az.BuildDefinition  `json:"definition"`
	Controller    *az.BuildController `json:"controller,omitempty"`
	LastChangedBy *az.IdentityRef     `json:"lastChangedBy,omitempty"`
	DeletedBy     *az.IdentityRef     `json:"deletedBy,omitempty"`
	BuildNumber   string              `json:"buildNumber,omitempty"`
	FinishTime    string              `json:"finishTime,omitempty"`
	Branch        string              `json:"sourceBranch"`
	Repository    az.Repository       `json:"repository"`
	Demands       []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"demands"`
	Logs *struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"logs,omitempty"`
	Project *struct {
		Abbreviation string `json:"abbreviation"`
		Description  string `json:"description"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		Revision     int    `json:"revision"`
		State        string `json:"state"`
		URL          string `json:"url"`
		Visibility   string `json:"visibility"`
	} `json:"project,omitempty"`
	Properties          map[string]string
	Priority            string `json:"priority,omitempty"`
	BuildNumberRevision int    `json:"buildNumberRevision,omitempty"`
	Deleted             *bool  `json:"deleted,omitempty"`
	DeletedDate         string `json:"deletedDate,omitempty"`
	DeletedReason       string `json:"deletedReason,omitempty"`
	ID                  int    `json:"id,omitempty"`
	KeepForever         bool   `json:"keepForever,omitempty"`
	ChangedDate         string `json:"lastChangedDate,omitempty"`
	Params              string `json:"parameters,omitempty"`
	Quality             string `json:"quality,omitempty"`
	Queue               struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Pool *struct {
			ID       int    `json:"id"`
			IsHosted bool   `json:"is_hosted"`
			Name     string `json:"name"`
		} `json:"pool,omitempty"`
	} `json:"queue"`
	QueueOptions      map[string]string `json:"queue_options"`
	QueuePosition     *int              `json:"queuePosition,omitempty"`
	QueueTime         string            `json:"queueTime,omitempty"`
	RetainedByRelease *bool             `json:"retainedByRelease,omitempty"`
	RequestedBy       *az.IdentityRef   `json:"requestedBy,omitempty"`
	RequestedFor      *az.IdentityRef   `json:"requestedFor,omitempty"`
	Version           string            `json:"sourceVersion,omitempty"`
	StartTime         string            `json:"startTime,omitempty"`
	Status            string            `json:"status,omitempty"`
	Result            string            `json:"result,omitempty"`
	ValidationResults []struct {
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	Tags         []string `json:"tags,omitempty"`
	TriggerBuild *Build   `json:"triggeredByBuild,omitempty"`
	URI          string   `json:"uri,omitempty"`
	URL          string   `json:"url,omitempty"`
}

func (a *AzureDevopsBuilds) Refresh() {
	a.getBuilds()
	return
}

// getBuilds fetches all builds under buildDefinitionIds
// on error, set err, and do not modify previously fetched builds
func (a *AzureDevopsBuilds) getBuilds() {
	if len(a.buildDefinitionIds) == 0 {
		return
	}

	params := url.Values{}
	params.Add("definitions", strings.Join(a.buildDefinitionIds, ","))
	params.Add("queryOrder", "queueTimeDescending")
	params.Add("maxBuildsPerDefinition", strconv.Itoa(a.numBuildsToFetch))

	// list builds
	URL := fmt.Sprintf(
		"_apis/build/builds?%s&%s",
		params.Encode(),
		"api-version=5.0",
	)

	request, err := a.client.NewRequest("GET", URL, nil)
	if err != nil {
		a.err = err
		return
	}

	var buildResp BuildsListResponse
	_, err = a.client.Execute(request, &buildResp)
	if err != nil {
		a.err = err
		return
	}

	// clear leftover builds
	a.builds = make(map[string][]Build)

	for _, build := range buildResp.Builds {
		sortKey := strconv.Itoa(build.Definition.ID)
		a.builds[sortKey] = append(a.builds[sortKey], build)
	}

	return
}

func constructClientFromConfig() *az.Client {
	return az.NewClient(
		wtf.Config.UString(
			"wtf.mods.azuredevopsrepo.account",
			os.Getenv("WTF_AZUREDEVOPS_ACCOUNT"),
		),
		wtf.Config.UString(
			"wtf.mods.azuredevopsrepo.project",
			os.Getenv("WTF_AZUREDEVOPS_PROJECT"),
		),
		wtf.Config.UString(
			"wtf.mods.azuredevopsrepo.token",
			os.Getenv("WTF_AZUREDEVOPS_TOKEN"),
		),
	)
}
