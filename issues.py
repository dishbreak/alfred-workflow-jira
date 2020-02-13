import sys
from workflow import Workflow, web, ICON_INFO
import base64
from config import get_missing_configs, JIRA_API_KEY, JIRA_URL, JIRA_USERNAME, JQL_QUERY
from utils import pluck_issue


def get_issues(url, user, password, query, wf):
	basic_auth = base64.b64encode("{}:{}".format(user, password))
	response = web.get(url + "/rest/api/2/search",
		params=dict(jql=query,),
		headers=dict(Authorization="Basic {}".format(basic_auth)),
	)

	wf.logger.debug(response.url)
	response.raise_for_status()
	data = response.json()
	wf.logger.debug(data)
	issues = [ pluck_issue(x, url) for x in data['issues'] ]
	max_issues = data['maxResults']
	total_issues = data['total']

	return issues, max_issues, total_issues


def main(wf):
	missing_configs = get_missing_configs(wf)
	if missing_configs:
		raise Exception("Missing configs: {}. Run `jirasetup` to configure".format(missing_configs))

	query = wf.args[0] if len(wf.args) > 0 else None

	issues, max_issues, total_issues = get_issues(
		wf.settings[JIRA_URL],
		wf.settings[JIRA_USERNAME],
		wf.get_password(JIRA_API_KEY),
		JQL_QUERY,
		wf)

	if query:
		issues = wf.filter(query, issues, key=lambda x: x['title'])

	wf.logger.debug(issues)
	if total_issues > max_issues:
		wf.add_item(title="{} issues found, showing the first {}".format(total_issues, max_issues),
			subtitle="Press enter to open search in Jira",
			arg="search", icon=ICON_INFO, valid=True)
	else:
		wf.add_item(title="Search in Jira", arg="search", icon=ICON_INFO, valid=True)

	for item in issues:
		wf.add_item(**item)

	wf.send_feedback()

if __name__ == '__main__':
	wf = Workflow()
	sys.exit(wf.run(main))
