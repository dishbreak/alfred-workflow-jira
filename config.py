from workflow import PasswordNotFound, Workflow
import sys

JIRA_API_KEY = "jira_api_key"
JIRA_URL = "jira_url"
JIRA_USERNAME = "jira_username"
JQL_QUERY = 'status in ("In Progress", "To Do", Triage) AND updated >= -52w AND assignee in (currentUser()) order by lastViewed DESC'

def get_missing_configs(wf):
	"""
	Try to verify that all configs are in place.
	"""
	missing_configs = []
	for key in [JIRA_URL, JIRA_USERNAME]:
		if key not in wf.settings:
			missing_configs.append(key)

	try:
		wf.get_password(JIRA_API_KEY)
	except PasswordNotFound:
		missing_configs.append(JIRA_API_KEY)

	return missing_configs

def parse_args(args):
	if len(args) != 2:
		raise Exception("Need exactly 2 arguments")

	parsed_args = {}

	if args[0] not in ['set_url', 'set_username', 'set_token', 'clear_config']:
		raise Exception("Unexpected command '{}'".format(args[0]))

	parsed_args['command'] = args[0]
	parsed_args['parameter'] = args[1]
	return parsed_args


def set_config(wf, key, value):
	wf.settings[key] = value
	wf.settings.save()

def set_password(wf, key, value):
	wf.save_password(key, value)

def clear_config(wf, value):
	if value:
		for key in [JIRA_URL, JIRA_USERNAME]:
			wf.settings.pop(key, None)
		wf.settings.save()
		try:
			wf.delete_password(JIRA_API_KEY)
		except PasswordNotFound:
			wf.logger.debug("No token found, completed.")

_dispatcher = {
	"set_url": lambda wf, value: set_config(wf, JIRA_URL, value),
	"set_username": lambda wf, value: set_config(wf, JIRA_USERNAME, value),
	"set_token": lambda wf, value: set_password(wf, JIRA_API_KEY, value),
	"clear_config": lambda wf, value: clear_config(wf, value)
}

def main(wf):
	args = parse_args(wf.args)
	wf.logger.debug(args)
	_dispatcher[args['command']](wf, args['parameter'])

if __name__ == '__main__':
	wf = Workflow()
	sys.exit(wf.run(main))
