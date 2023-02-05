package main

import (
	"fmt"

	"github.com/dishbreak/go-alfred/alfred"
)

const (
	JiraUsername = "jira_username"
	JiraUrl      = "jira_url"
	JiraBoard    = "jira_board"
	JqlQuery     = `status in ("In Progress", "To Do", Triage, "Code Review") AND updated >= -52w AND assignee in (currentUser()) order by lastViewed DESC`
)

type SetConfigCmd struct {
	ConfigKey   string `arg:"" required:""`
	ConfigValue string `arg:"" required:""`
}

func (s *SetConfigCmd) Run(ctx *Context) error {
	alfred.RunScriptAction(func(sar *alfred.ScriptActionResponse) error {
		switch s.ConfigKey {
		case JiraUsername:
		case JiraBoard:
		case JiraUrl:
			return ctx.wf.SetConfig(s.ConfigKey, s.ConfigValue)
		default:
			return fmt.Errorf("unrecognized config '%s'", s.ConfigKey)
		}
		return nil
	})
	return nil
}
