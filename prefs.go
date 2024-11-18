package main

import (
	"io"

	"github.com/groob/plist"
)

type Prefs struct {
	MyBoardId string `plist:"jira_board_id"`
	JiraEmail string `plist:"jira_email"`
	JiraUrl   string `plist:"jira_url"`
}

func loadPrefs(r io.Reader) (*Prefs, error) {
	dec := plist.NewXMLDecoder(r)
	p := &Prefs{}
	err := dec.Decode(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
