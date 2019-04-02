package azuredevopsrepo

import (
	"fmt"
	"net/url"
	"os"

	az "github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/wtfutil/wtf/wtf"
)

type AzureDevopsRepo struct {
	client       *az.Client
	Repo         az.Repository
	PullRequests []PullRequest

	// if an operation resulted in an error, it should be stored here
	// so that it can be displayed
	err error
}

type PullRequestsResponse struct {
	PullRequests []PullRequest `json:"value"`
	Count        int           `json:"count"`
}

type PullRequest struct {
	ID          int                `json:"pullRequestId,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Created     string             `json:"creationDate"`
	CreatedBy   User               `json:"createdBy"`
	Repo        az.PullRequestRepo `json:"repository"`
	URL         string             `json:"url"`
	Reviewers   []User             `json:"reviewers"`
}

type User struct {
	Vote          int    `json:"vote,omitempty"`
	ID            string `json:"id"`
	DisplayName   string `json:"displayName"`
	UniqueName    string `json:"uniqueName"`
	IsAadIdentity bool   `json:"isAadIdentity"`
	IsContainer   bool   `json:"isContainer"`
}

func NewRepo(repoName string) (r *AzureDevopsRepo) {
	r = &AzureDevopsRepo{}
	r.client = constructClientFromConfig()

	URL := fmt.Sprintf(
		"/_apis/git/repositories/%s?api-version=4.1",
		url.PathEscape(repoName),
	)

	var azrepo az.Repository
	request, err := r.client.NewRequest("GET", URL, nil)
	if err != nil {
		r.err = err
		return
	}

	_, err = r.client.Execute(request, &azrepo)
	if err != nil {
		r.err = err
	}
	r.Repo = azrepo

	return
}

func (r *AzureDevopsRepo) Refresh() {
	var errs []error
	if err := r.loadPullRequests(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		r.err = fmt.Errorf("Error(s) occurred: %v", errs)
	}

	return
}

func (r *AzureDevopsRepo) loadPullRequests() error {
	params := url.Values{}
	params.Add("searchCriteria.repositoryId", r.Repo.ID)

	URL := fmt.Sprintf(
		"/_apis/git/pullrequests?%s&%s",
		"api-version=4.1",
		params.Encode(),
	)

	request, err := r.client.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	var response PullRequestsResponse
	_, err = r.client.Execute(request, &response)
	if err != nil {
		return err
	}

	r.PullRequests = response.PullRequests
	return nil
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

func containsUser(name string, Users ...User) bool {
	for _, user := range Users {
		if user.DisplayName == name || user.UniqueName == name {
			return true
		}
	}
	return false
}
