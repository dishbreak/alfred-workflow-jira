from client.base import BaseClient


class AgileClient(BaseClient):
    def get_boards(self):
       return self.get_from_api_paging("rest/agile/1.0/board", params={})

    def get_issues_for_board(self, board_id, jql=None, fields=list()):
        params = dict()
        if jql:
            params['jql'] = jql
        if fields:
            params['fields'] = ",".join(fields)

        route = "rest/agile/1.0/board/%s/issue" % board_id
        return self.get_from_api_paging(route, params, page_key="issues")

    def get_configuration_for_board(self, board_id):
        route = "rest/agile/1.0/board/%s/configuration" % board_id
        resp = self.get_from_api(route, params={})
        return resp.json()

