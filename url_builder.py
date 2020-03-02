from workflow import Workflow
from workflow.notify import notify
import sys
from config import JIRA_URL, get_missing_configs, JQL_QUERY
import urllib
from urlparse import urljoin


def parse_args(argv):
    args = {}
    if argv[0] not in ['issue', 'search', 'board']:
        raise Exception('unexpected argument {}'.format(argv[0]))

    args['url_type'] = argv[0]

    if args['url_type'] == "issue":
        if len(argv) != 2:
            raise Exception("issue url needs a second argument")
        args['issue_key'] = argv[1].strip()

    elif args['url_type'] == "board":
        if len(argv) != 2:
            raise Exception("board url needs a second argument")
        args['board_id'] = argv[1].strip()

    return args

def get_issue_url(url, issue_key):
    return urljoin(url, "browse/%s" % issue_key)

def get_search_url(url, jql):
    encoded = urllib.urlencode(dict(jql=jql))
    return urljoin(url, "/issues/?%s" % encoded)

def get_board_url(url, board_id):
    return urljoin(url, "secure/RapidBoard.jspa?rapidView=%s" % board_id)

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
    elif args['url_type'] == "board":
        result = get_board_url(
            wf.settings[JIRA_URL],
            args['board_id'],
        )
    wf.logger.debug(result)
    print(result)

if __name__ == '__main__':
    wf = Workflow()
    sys.exit(wf.run(main))
