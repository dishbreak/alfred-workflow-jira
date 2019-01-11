import sys
from workflow import Workflow, web, ICON_INFO
import base64

API_TOKEN = 'bEMnM7sEGEyTYq0YyMiY9C70'
USER = 'vkotcherlakota@nerdwallet.com'
URL = "https://nerdwallet.atlassian.net"
JQL_QUERY = 'status in ("In Progress", "To Do", Triage) AND updated >= -52w AND assignee in (currentUser()) order by lastViewed DESC'

def pluck_issue(issue):
	plucked = {"valid": True}
	plucked['arg'] = issue['key']
	plucked['title'] = issue['key'] + ": " + issue['fields']['summary']
	plucked['subtitle'] = issue['fields']['description'][0:80] + "..." if len(issue['fields']['description']) > 80 else "" 
	return plucked

def get_issues(url, user, password, wf):
	basic_auth = base64.b64encode("{}:{}".format(user, password))
	response = web.get(url + "/rest/api/2/search",
		params=dict(jql=JQL_QUERY,),
		headers=dict(Authorization="Basic {}".format(basic_auth)),
	)

	wf.logger.debug(response.url)
	response.raise_for_status()
	data = response.json()
	issues = [ pluck_issue(x) for x in data['issues'] ]
	max_issues = data['maxResults']
	total_issues = data['total']

	return issues, max_issues, total_issues


def main(wf):
	query = wf.args[0] if len(wf.args) > 0 else None

	issues, max_issues, total_issues = get_issues(URL, USER, API_TOKEN, wf)
	if query:
		issues = wf.filter(query, issues, key=lambda x: x['title'])

	wf.logger.debug(issues)
	if total_issues > max_issues:
		wf.add_item(title="{} issues found, showing the first {}".format(total_issues, max_issues),
			subtitle="Press enter to open search in Jira",
			arg="search", icon=ICON_INFO)

	for item in issues:
		wf.add_item(**item)

	wf.send_feedback()

if __name__ == '__main__':
	wf = Workflow()
	sys.exit(wf.run(main))
