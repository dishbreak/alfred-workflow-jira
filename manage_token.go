package main

import (
	"github.com/dishbreak/alfred-workflow-jira/keychain"
)

type SaveTokenCmd struct {
	TokenValue string `arg:"" help:"The token value to save."`
}

func (cmd *SaveTokenCmd) Run(c *Context) error {
	return keychain.SaveToken(cmd.TokenValue)
}

type DeleteTokenCmd struct{}

func (cmd *DeleteTokenCmd) Run(c *Context) error {
	return keychain.DeleteToken()
}
