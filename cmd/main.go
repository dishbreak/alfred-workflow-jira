package main

import (
	"github.com/alecthomas/kong"
	"github.com/andygrunwald/go-jira"
	"github.com/dishbreak/go-alfred/alfred"
)

type Context struct {
	wf *alfred.Workflow
}

type CtxErr string

func (c CtxErr) Error() string {
	return string(c)
}

var cli struct {
	Setup struct {
		SetToken    SetTokenCmd    `cmd:"set-token" help:"save an API token"`
		ForgetToken ForgetTokenCmd `cmd:"forget-token" help:"clear an API token"`
		SetConfig   SetConfigCmd   `cmd:"set-config" help:"set a config value"`
	} `cmd:"setup" help:"commands used to configure workflow"`
}

const (
	ErrConfigNotSet CtxErr = "missing configuration"
)

func (ctx *Context) GetJiraClient() (*jira.Client, error) {
	userName := ctx.wf.GetConfigString(JiraUsername, "unset")
	if userName == "unset" {
		return nil, ErrConfigNotSet
	}

	url := ctx.wf.GetConfigString(JiraUrl, "unset")
	if userName == "unset" {
		return nil, ErrConfigNotSet
	}

	token, err := ctx.wf.GetSecret(JiraApiToken)
	if err != nil {
		return nil, err
	}

	tp := jira.BasicAuthTransport{
		Username: userName,
		Password: token,
	}

	return jira.NewClient(tp.Client(), url)
}

const (
	workflowName = "jira-browser"
)

func main() {
	wf, err := alfred.NewWorkflow(workflowName)
	if err != nil {
		panic(err)
	}

	ctx := kong.Parse(&cli)
	err = ctx.Run(&Context{
		wf: wf,
	})
	ctx.FatalIfErrorf(err)
}
