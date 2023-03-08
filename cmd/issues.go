package main

import (
	"fmt"
)

type ListIssuesCmd struct {
	SubstringQuery string `arg:"" default:""`
}

func (l *ListIssuesCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		jiraClient, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		issues, resp, err := jiraClient.Issue.Search(JqlQuery, nil)
		if err != nil {
			panic(err)
		}
		searchTitle := "Search in Jira"
		if resp.Total > resp.MaxResults {
			searchTitle = fmt.Sprintf("%d results found, showing first %d", resp.Total, resp.MaxResults)
		}

		ctx.wf.NewItem(searchTitle).Subtitle("Open in Jira").Valid(true).Arg("search")

		for _, issue := range issues {
			ctx.RenderIssue(&issue)
		}

		if l.SubstringQuery != "" {
			ctx.wf.Filter(l.SubstringQuery)
		}

		ctx.wf.SendFeedback()
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
