from base_runner import BaseRunner
from client.agile import AgileClient
import config
from workflow import Workflow, ICON_GROUP, ICON_WEB
from collections import defaultdict
import sys
import itertools
from utils import pluck_issue
from urlparse import urljoin


JQL_QUERY = "sprint in openSprints() and sprint not in futureSprints() and assignee in (currentUser())"

class BoardRunner(BaseRunner):

    def main(self, wf):
        client = AgileClient(
            url=wf.settings[config.JIRA_URL],
            user=wf.settings[config.JIRA_USERNAME],
            token=wf.get_password(config.JIRA_API_KEY),
            logger=wf.logger,
        )

        selected_board = wf.args[1] if len(wf.args) > 1 else None

        if not selected_board and config.JIRA_BOARD_ID not in wf.settings:
            raise Exception("Missing Jira Board ID. Run jiraconfig and select a board.")

        selected_board = selected_board if selected_board else wf.settings[config.JIRA_BOARD_ID]

        issues = client.get_issues_for_board(
            board_id=selected_board,
            jql=JQL_QUERY,
            fields=["status", "key", "description", "summary"],
        )

        conf = client.get_configuration_for_board(
            board_id=selected_board
        )

        id_to_category = {}
        column_names = []

        for column in conf['columnConfig']['columns']:
            column_names.append(column['name'])
            for status in column['statuses']:
                id_to_category[status['id']] = column['name']

        wf.logger.debug("id_to_category: %s", id_to_category)

        sorted(issues, key=lambda x: x['fields']['status']['id'])
        def keyfunc(x):
            return id_to_category[x['fields']['status']['id']]

        issues_by_column = defaultdict(list)
        for col, issues_in_col in itertools.groupby(issues, key=keyfunc):
            wf.logger.debug("col: %s issues %s", col, issues_in_col)
            issues_by_column[col] += list(issues_in_col)

        board_jira_url = urljoin(
            wf.settings[config.JIRA_URL],
            "secure/RapidBoard.jspa?rapidView=%s" % selected_board
        )

        wf.add_item(title="View In Jira", icon=ICON_WEB, arg=board_jira_url, valid=True)
        for column in column_names:
            wf.add_item(title=column, icon=ICON_GROUP)
            for issue in issues_by_column[column]:
                wf.add_item(**pluck_issue(
                    issue,
                    wf.settings[config.JIRA_URL]
                ))

        wf.send_feedback()

class BoardPickerRunner(BaseRunner):
    def main(self, wf):
        client = AgileClient(
            url=wf.settings[config.JIRA_URL],
            user=wf.settings[config.JIRA_USERNAME],
            token=wf.get_password(config.JIRA_API_KEY),
            logger=wf.logger,
        )

        boards = client.get_boards()
        sorted(boards, key=lambda x: x["name"])

        for board in boards:
            wf.add_item(title=board["name"], arg=str(board["id"]), valid=True)

        wf.send_feedback()


def main(wf):
    wf.logger.debug("got to loop")
    command = wf.args[0] if len(wf.args) > 0 else None
    runner = None
    wf.logger.debug("command: %s", command)
    if command == "get_board_issues":
        runner = BoardRunner()
    elif command == "get_boards":
        runner = BoardPickerRunner()
    else:
        raise NotImplementedError("Don't support command '%s'" % command)

    runner.run(wf)


if __name__ == '__main__':
    wf = Workflow()
    wf.logger.debug("hi")
    sys.exit(wf.run(main))
