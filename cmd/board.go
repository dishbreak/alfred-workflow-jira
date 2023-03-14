package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/andygrunwald/go-jira"
	aw "github.com/deanishe/awgo"
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
	BoardID      int `default:"-1"`
	OnlyMyIssues bool
}

type issuesForBoardResp struct {
	StartAt    int          `json:"startAt"`
	MaxResults int          `json:"maxResults"`
	Total      int          `json:"total"`
	Issues     []jira.Issue `json:"issues"`
}

const (
	BoardJql     = "Sprint in openSprints() AND Sprint not in futureSprints()"
	OnlyMyIssues = " AND assignee=currentUser()"
)

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
			panic(errors.New("no board id set"))
		}

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

		jql := BoardJql
		if l.OnlyMyIssues {
			jql += OnlyMyIssues
		}

		for startAt := 0; ; {
			r, err := j.NewRequest("GET", fmt.Sprintf("/rest/agile/1.0/board/%d/issue", l.BoardID), nil)
			if err != nil {
				panic(fmt.Errorf("failed to form api request: %w", err))
			}
			q := r.URL.Query()
			q.Add("startAt", strconv.Itoa(startAt))
			q.Add("jql", jql)
			r.URL.RawQuery = q.Encode()

			var result issuesForBoardResp
			_, err = j.Do(r, &result)
			if err != nil {
				panic(fmt.Errorf("failed to get issues for board: %w", err))
			}

			issues = append(issues, result.Issues...)
			startAt += len(result.Issues)
			if startAt >= result.Total {
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
			ctx.wf.NewItem(column.Name).Icon(aw.IconGroup)
			for _, issue := range issuesByCol[column.Name] {
				ctx.RenderIssue(&issue)
			}
		}

		ctx.wf.SendFeedback()
	})
	return nil
}
