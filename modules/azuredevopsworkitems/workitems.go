package azuredevopsworkitems

import (
	"fmt"
	"os"
	"strconv"

	az "github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/wtfutil/wtf/wtf"
)

type AzureDevopsWorkItems struct {
	client      *az.Client
	queryId     string
	workItemIds []int
	workItems   []WorkItem

	err error
}

type QueryResp struct {
	WorkItems []WorkItemReference `json:"workItems"`
}

type WorkItemReference struct {
	ID int `json:"id"`
}

type WorkItemResp struct {
	Count     int        `json:"count"`
	WorkItems []WorkItem `json:"value"`
}

type WorkItem struct {
	ID     int            `json:"id"`
	Fields WorkItemFields `json:"fields"`
}

type WorkItemFields struct {
	Title        string `json:"System.Title"`
	Description  string `json:"System.Description"`
	State        string `json:"System.State"`
	WorkItemType string `json:"System.WorkItemType"`
	AssignedTo   User   `json:"System.AssignedTo"`
}

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}

func (a *AzureDevopsWorkItems) Refresh() {
	a.loadWorkItemIdsFromQuery()
	a.loadWorkItemDetails()
	return
}

// on failure, do not modify current list of work item ids
func (a *AzureDevopsWorkItems) loadWorkItemIdsFromQuery() {
	URL := fmt.Sprintf(
		"/_apis/wit/wiql/%s?%s",
		a.queryId,
		"api-version=5.0",
	)

	request, err := a.client.NewRequest("GET", URL, nil)
	if err != nil {
		a.err = err
		return
	}

	var resp QueryResp
	_, err = a.client.Execute(request, &resp)
	if err != nil {
		a.err = err
		return
	}

	ids := []int{}
	for _, ref := range resp.WorkItems {
		ids = append(ids, ref.ID)
	}

	a.workItemIds = ids
	return
}

func (a *AzureDevopsWorkItems) loadWorkItemDetails() {
	if len(a.workItemIds) == 0 {
		return
	}

	commaSeparatedIds := ""
	for _, id := range a.workItemIds {
		if commaSeparatedIds != "" {
			commaSeparatedIds = commaSeparatedIds + ","
		}
		commaSeparatedIds = commaSeparatedIds + strconv.Itoa(id)
	}

	URL := fmt.Sprintf(
		"_apis/wit/workitems?ids=%s&$expand=all&%s",
		commaSeparatedIds,
		"api-version=5.0",
	)

	request, err := a.client.NewRequest("GET", URL, nil)
	if err != nil {
		a.err = err
		return
	}

	var resp WorkItemResp
	_, err = a.client.Execute(request, &resp)
	if err != nil {
		a.err = err
		return
	}

	a.workItems = resp.WorkItems
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
