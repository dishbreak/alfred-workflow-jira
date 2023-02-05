package main

import (
	"fmt"
	"net/url"

	"github.com/dishbreak/go-alfred/alfred"
)

type ListIssuesCmd struct {
	SubstringQuery string `arg:"" default:""`
}

func (l *ListIssuesCmd) Run(ctx *Context) error {
	alfred.RunScriptFilter(func(sfr *alfred.ScriptFilterResponse) error {
		jiraClient, err := ctx.GetJiraClient()
		if err != nil {
			return err
		}

		baseUrl := ctx.wf.GetConfigString(JiraUrl, "")
		u, err := url.Parse(baseUrl)
		if err != nil {
			return err
		}

		issues, resp, err := jiraClient.Issue.Search(JqlQuery, nil)
		if err != nil {
			return err
		}
		searchTitle := "Search in Jira"
		if resp.Total > resp.MaxResults {
			searchTitle = fmt.Sprintf("%d results found, showing first %d", resp.Total, resp.MaxResults)
		}

		sfr.AddItem(alfred.ListItem{
			Title:    searchTitle,
			Subtitle: "Open in Jira",
			Valid:    true,
			Arg:      "search",
		})

		for _, issue := range issues {
			issueUrl, _ := u.Parse("browse/" + issue.Key)
			sfr.AddItem(alfred.ListItem{
				Title:    fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary),
				Subtitle: truncate(issue.Fields.Description, 89, "No description given."),
				Arg:      issueUrl.String(),
				Valid:    true,
			})
		}
		return nil
	})
	return nil
}

func truncate(input string, maxLen int, defaultText string) string {
	if input == "" {
		return defaultText
	}
	if len(input) > maxLen {
		return input[:maxLen-3] + "..."
	}
	return input
}
