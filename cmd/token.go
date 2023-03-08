package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type SetTokenCmd struct{}

const (
	JiraApiToken = "jira_api_token"
)

func (s *SetTokenCmd) Run(ctx *Context) error {

	ctx.wf.Run(func() {
		var token string
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		token = string(b)

		token = strings.TrimSpace(token)

		err = ctx.wf.Keychain.Set(JiraApiToken, token)
		if err != nil {
			panic(err)
		}
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
	ctx.wf.Run(func() {
		if err := ctx.wf.Keychain.Delete(JiraApiToken); err != nil {
			panic(err)
		}
	})
	return nil
}

func (f *ForgetTokenCmd) Help() string {
	return "Removes the Atlassian Token from the system keychain."
}
