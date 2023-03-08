package main

import (
	"fmt"
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/andygrunwald/go-jira"
	aw "github.com/deanishe/awgo"
)

type Context struct {
	wf      *aw.Workflow
	baseUrl *url.URL
}

type CtxErr string

func (c CtxErr) Error() string {
	return string(c)
}

var cli struct {
	Setup struct {
		SetToken    SetTokenCmd    `cmd:"set-token" help:"save an API token"`
		ForgetToken ForgetTokenCmd `cmd:"forget-token" help:"clear an API token"`
	} `cmd:"setup" help:"commands used to configure workflow"`
	Board struct {
		ListBoards ListBoardsCmd         `cmd:"list-boards" help:"list all matching boards"`
		ListIssues ListIssuesForBoardCmd `cmd:"list-issues" help:"list all issues on the given board"`
	} `cmd:"board" help:"commands used to interact with boards"`
	Issue struct {
		List ListIssuesCmd `cmd:"list" help:"list issues matching the given query"`
	} `cmd:"issue" help:"commands used to interact with issues"`
}

const (
	ErrConfigNotSet CtxErr = "missing configuration"
)

func (ctx *Context) GetJiraClient() (*jira.Client, error) {
	userName := ctx.wf.Config.Get(JiraUsername, "unset")
	if userName == "unset" {
		return nil, ErrConfigNotSet
	}

	url := ctx.wf.Config.Get(JiraUrl, "unset")
	if userName == "unset" {
		return nil, ErrConfigNotSet
	}

	token, err := ctx.wf.Keychain.Get(JiraApiToken)
	if err != nil {
		return nil, err
	}

	tp := jira.BasicAuthTransport{
		Username: userName,
		Password: token,
	}

	return jira.NewClient(tp.Client(), url)
}

func (c *Context) RenderIssue(issue *jira.Issue) {
	issueUrl, _ := c.baseUrl.Parse("browse/" + issue.Key)
	item := c.wf.NewItem(fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary))
	item.Subtitle(truncate(issue.Fields.Description, 80, "No description given."))
	item.Arg(issueUrl.String()).Valid(true)
}

const (
	workflowName = "jira-browser"
)

func main() {
	wf := aw.New()

	baseUrl := wf.Config.Get(JiraUrl, "")
	u, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}

	ctx := kong.Parse(&cli)
	err = ctx.Run(&Context{
		wf:      wf,
		baseUrl: u,
	})
	ctx.FatalIfErrorf(err)
}
