package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

type ListFiltersCmd struct{}

func (l *ListFiltersCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		j, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		filters, _, err := j.Filter.GetFavouriteList()
		if err != nil {
			panic(err)
		}

		for _, filter := range filters {
			ctx.wf.NewItem(filter.Name).Arg(filter.ID).Valid(true)
		}

	})

	return nil
}

type QueryFilterCmd struct {
	FilterID string `arg:""`
}

func (q *QueryFilterCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		j, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		jql := fmt.Sprintf("filter = %s", q.FilterID)

		for startAt := 0; ; {
			chunk, resp, err := j.Issue.Search(jql, &jira.SearchOptions{
				StartAt: startAt,
			})
			if err != nil {
				panic(err)
			}

			for _, issue := range chunk {
				ctx.RenderIssue(&issue)
			}
			startAt += resp.Total
			if startAt > resp.Total {
				break
			}
		}
		ctx.wf.SendFeedback()
	})

	return nil
}
