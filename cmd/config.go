package main

const (
	JiraUsername = "jira_username"
	JiraUrl      = "jira_url"
	JiraBoard    = "jira_board"
	JqlQuery     = `status in ("In Progress", "To Do", Triage, "Code Review") AND updated >= -52w AND assignee in (currentUser()) order by lastViewed DESC`
)
