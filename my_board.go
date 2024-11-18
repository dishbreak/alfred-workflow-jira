package main

import (
	"fmt"

	"github.com/dishbreak/alfred-workflow-jira/jira"
	"github.com/dishbreak/go-alfred/v2/alfred"
)

type MyBoardCmd struct{}

func (m *MyBoardCmd) Run(c *Context) error {
	sfr := alfred.NewScriptFilterResponse()

	defer alfred.RecoverIfErr(sfr)()

	if c.jiraClient == nil {
		panic("Missing Jira credentials. Run jirasetup to set token")
	}

	bc, err := c.jiraClient.GetBoardConfig(c.prefs.MyBoardId)
	if err != nil {
		panic(fmt.Errorf("failed to fetch board config for board %s: %w", c.prefs.MyBoardId, err))
	}

	sprintList, err := c.jiraClient.GetBoardActiveSprints(c.prefs.MyBoardId)
	if err != nil {
		panic(fmt.Errorf("failed to get board sprints: %w", err))
	}

	var issues []jira.Issue
	for _, sprint := range sprintList.Values {
		result, err := c.jiraClient.Search(
			fmt.Sprintf(`sprint = %d AND assignee = currentUser()`, sprint.Id),
			jira.WithSearchFields([]string{
				"description",
				"status",
				"summary",
			}),
		)
		if err != nil {
			fmt.Printf("failed to search for issues in sprint '%s': %s\n", sprint.Name, err)
			continue
		}
		issues = append(issues, result.Issues...)
	}

	issuesByStatusID := make(map[string][]jira.Issue)
	for _, issue := range issues {
		status := issue.Fields.Status.ID
		slice := issuesByStatusID[status]
		slice = append(slice, issue)
		issuesByStatusID[status] = slice
	}

	for _, col := range bc.ColumnConfig.Columns {
		sfr.AddItem(alfred.ListItem{
			Title: col.Name,
			Valid: false,
		})
		for _, status := range col.Statuses {
			for _, issue := range issuesByStatusID[status.Id] {
				sfr.AddItem(alfred.ListItem{
					Title:    fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary),
					Subtitle: issue.Fields.Description,
					Valid:    true,
					Arg:      issue.Key,
				})
			}
		}
	}

	sfr.SendFeedback()

	return nil
}
