package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/andygrunwald/go-jira"
)

type ListBoardsCmd struct{}

func (l *ListBoardsCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		j, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		boards, _, err := j.Board.GetAllBoards(&jira.BoardListOptions{})
		if err != nil {
			panic(err)
		}

		for _, board := range boards.Values {
			it := ctx.wf.NewItem(board.Name)
			it.Arg(strconv.Itoa(board.ID)).Valid(true)
		}

		ctx.wf.SendFeedback()
	})

	return nil
}

type ListIssuesForBoardCmd struct {
	BoardID int `arg:"" required:""`
}

func (l *ListIssuesForBoardCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		j, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		b, _, err := j.Board.GetBoard(l.BoardID)
		if err != nil {
			panic(err)
		}

		query := fmt.Sprintf("filter = %d AND sprint in openSprints() and sprint not in futureSprints() and assignee in (currentUser())", b.FilterID)

		conf, _, err := j.Board.GetBoardConfiguration(l.BoardID)
		if err != nil {
			panic(err)
		}

		stateMap := make(map[string]string)
		issuesByCol := make(map[string][]jira.Issue)

		for _, col := range conf.ColumnConfig.Columns {
			for _, status := range col.Status {
				stateMap[status.ID] = col.Name
			}
			issuesByCol[col.Name] = make([]jira.Issue, 0)
		}

		issues := make([]jira.Issue, 0)
		paginate := func(startAt int) *jira.SearchOptions {
			return &jira.SearchOptions{
				StartAt: startAt,
			}
		}

		for startAt := 0; ; {
			list, resp, err := j.Issue.Search(query, paginate(startAt))
			if err != nil {
				panic(err)
			}

			issues = append(issues, list...)
			startAt += resp.Total
			if startAt > resp.MaxResults {
				break
			}
		}

		for _, issue := range issues {
			colName, ok := stateMap[issue.Fields.Status.ID]
			if !ok {
				log.Printf("issue %s has status %s, no matching column in board %d",
					issue.Key, issue.Fields.Status.Name, l.BoardID)
				continue
			}
			issuesByCol[colName] = append(issuesByCol[colName], issue)
		}

		for _, column := range conf.ColumnConfig.Columns {
			ctx.wf.NewItem(column.Name)
			for _, issue := range issuesByCol[column.Name] {
				ctx.RenderIssue(&issue)
			}
		}

		ctx.wf.SendFeedback()
	})
	return nil
}
