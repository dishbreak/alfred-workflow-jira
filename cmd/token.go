package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dishbreak/go-alfred/alfred"
)

type SetTokenCmd struct{}

const (
	JiraApiToken = "jira_api_token"
)

func (s *SetTokenCmd) Run(ctx *Context) error {

	alfred.RunScriptAction(func(sar *alfred.ScriptActionResponse) error {
		var token string
		if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			token = string(b)
		} else {
			fmt.Println("Enter user token. Press <ENTER> when done.")
			reader := bufio.NewReader(os.Stdin)
			token, _ = reader.ReadString('\n')
		}

		token = strings.TrimSpace(token)

		return ctx.wf.SaveSecret(JiraApiToken, token)
	})
	return nil
}

func (s *SetTokenCmd) Help() string {
	return `
Set an API for access to Atlassian APIs. For details, check documentation:
https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/

Note that for the sake of security, tokens must be passed via stdin. For example:
	echo "TOKEN" | jira-browser setup token
When run outside of a pipe, the command will prompt you for a token on input.
`
}

type ForgetTokenCmd struct{}

func (f *ForgetTokenCmd) Run(ctx *Context) error {
	alfred.RunScriptAction(func(sar *alfred.ScriptActionResponse) error {
		return ctx.wf.DeleteSecret(JiraApiToken)
	})

	return nil
}

func (f *ForgetTokenCmd) Help() string {
	return "Removes the Atlassian Token from the system keychain."
}
