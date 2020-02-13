from workflow import web
from urlparse import urljoin
import base64

class PagingException(Exception):
    """
    Raised when one tries using BaseClient.get_from_api_paged() for an endpoint
    that does not use paging.
    """
    pass

class BaseClient(object):

    def __init__(self, url, user, token, logger):
        self.url = url
        self.user = user
        self.token = token
        self.logger = logger
        self.basic_auth_header = "Basic {}".format(base64.b64encode("{}:{}".format(user, token)))

    def to_url(self, route):
        return urljoin(self.url, route)

    def get_from_api(self, route, params):
        url = self.to_url(route)
        self.logger.info(
            "api_request %s (%s) - %s",
            url,
            route,
            params,
        )

        resp = web.get(
            url,
            params=params,
            headers=dict(Authorization=self.basic_auth_header)
        )

        resp.raise_for_status()
        self.logger.debug(resp.text)

        return resp

    def get_from_api_paging(self, route, params, page_key="values"):
        keep_paging = True
        items = []
        params['startAt'] = params['startAt'] if 'startAt' in params else 0

        while keep_paging:

            resp = self.get_from_api(route, params)

            data = resp.json()
            params['startAt'] += data['maxResults']
            items += data[page_key]
            if 'isLast' in data:
                keep_paging = not data['isLast']
            elif 'total' in data:
                keep_paging = data['total'] > params
            else:
                raise Exception("can't page this api call")

        return items

