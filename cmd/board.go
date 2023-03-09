package main

import (
	"errors"
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

		for startAt := 0; ; {
			boards, _, err := j.Board.GetAllBoards(&jira.BoardListOptions{SearchOptions: jira.SearchOptions{StartAt: startAt}})
			if err != nil {
				panic(err)
			}

			for _, board := range boards.Values {
				it := ctx.wf.NewItem(board.Name)
				it.Arg(strconv.Itoa(board.ID)).Valid(true)
			}

			startAt += len(boards.Values)
			if startAt >= boards.Total {
				break
			}
		}

		ctx.wf.SendFeedback()
	})

	return nil
}

type SaveFavoriteBoardCmd struct {
	BoardID int `arg:"" required:""`
}

func (s *SaveFavoriteBoardCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		err := ctx.wf.Config.Set(JiraBoard, fmt.Sprintf("%d", s.BoardID), true).Do()
		if err != nil {
			panic(err)
		}

	})
	return nil
}

type ListIssuesForBoardCmd struct {
	BoardID int `default:"-1"`
}

func (l *ListIssuesForBoardCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		j, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		if l.BoardID == -1 {
			l.BoardID = ctx.wf.Config.GetInt(JiraBoard, -1)
		}

		if l.BoardID == -1 {
			panic(errors.New("no board specified"))
		}

		b, _, err := j.Board.GetBoard(l.BoardID)
		if err != nil {
			panic(fmt.Errorf("failed to get board: %w", err))
		}

		log.Println(b)

		query := fmt.Sprintf("filter = %d AND sprint in openSprints() and sprint not in futureSprints() and assignee in (currentUser())", b.FilterID)

		log.Printf("using query: %s", query)

		conf, _, err := j.Board.GetBoardConfiguration(l.BoardID)
		if err != nil {
			panic(fmt.Errorf("failed to get board config: %w", err))
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

		log.Printf("found %d issue(s) in board %d", len(issues), l.BoardID)

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
