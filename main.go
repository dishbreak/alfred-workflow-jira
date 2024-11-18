package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/dishbreak/alfred-workflow-jira/jira"
	"github.com/dishbreak/alfred-workflow-jira/keychain"
)

type CLI struct {
	SaveToken   SaveTokenCmd   `cmd:"" help:"Save token to Keychain"`
	DeleteToken DeleteTokenCmd `cmd:"" help:"Delete token from Keychain"`
	MyBoard     MyBoardCmd     `cmd:"" help:"Retrieve board contents"`
}

type Context struct {
	jiraClient *jira.Client
	prefs      *Prefs
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli)
	c := &Context{}

	f, err := os.Open("./prefs.plist")
	if err != nil {
		panic(fmt.Errorf("failed to load preferences: %w", err))
	}

	if p, err := loadPrefs(f); err != nil {
		panic(fmt.Errorf("failed to parse preferences: %w", err))
	} else {
		c.prefs = p
	}

	if token, err := keychain.GetToken(); err == nil {
		c.jiraClient = jira.NewClient(nil,
			jira.WithApiToken(token),
			jira.WithBaseUrl(c.prefs.JiraUrl),
			jira.WithUsername(c.prefs.JiraEmail),
		)
	} else {
		println("WARNING: failed to retrieve token from keychain: %s", err)
	}

	if err := ctx.Run(c); err != nil {
		fmt.Println(err)
		return
	}
}
