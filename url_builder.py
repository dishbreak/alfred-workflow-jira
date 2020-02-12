from workflow import Workflow
from workflow.notify import notify
import sys
from config import JIRA_URL, get_missing_configs, JQL_QUERY
import urllib

def parse_args(argv):
	args = {}
	if argv[0] not in ['issue', 'search']:
		raise Exception('unexpected argument {}'.format(argv[0]))

	args['url_type'] = argv[0]

	if args['url_type'] == "issue":
		if len(argv) != 2:
			raise Exception("issue url needs a second argument")
		args['issue_key'] = argv[1].strip()

	return args

def get_issue_url(url, issue_key):
	return url + "/browse/" + issue_key

def get_search_url(url, jql):
	encoded = urllib.urlencode(dict(jql=jql))
	return url + "/issues/?" + encoded

def main(wf):
	missing_configs = get_missing_configs(wf)
	if missing_configs:
		notify("Missing Configuration",
			"Run the `jirasetup` keyword to set config.")
		raise Exception("Missing configs: {}. Run `jirasetup` to configure".format(missing_configs))

	args = parse_args(wf.args)
	wf.logger.debug(args)
	result = None
	if args['url_type'] == "issue":
		result = get_issue_url(
			wf.settings[JIRA_URL],
			args['issue_key']
		)
	elif args['url_type'] == "search":
		result = get_search_url(
			wf.settings[JIRA_URL],
			JQL_QUERY
		)
	wf.logger.debug(result)
	print(result)

if __name__ == '__main__':
	wf = Workflow()
	sys.exit(wf.run(main))
