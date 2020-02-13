from urlparse import urljoin
import logging

logger = logging.getLogger(__name__)

def pluck_issue(issue, url):
    plucked = {"valid": True}
    logger.debug("issue: %s", issue)
    plucked['arg'] = urljoin(url, "browse/%s" % issue['key'])
    plucked['title'] = issue['key'] + ": " + issue['fields']['summary']
    plucked['subtitle'] = truncate(issue['fields']['description'], length=80, default_text="No description given.")
    return plucked

def truncate(string, length, default_text):
    if not string:
        return default_text
    return string[0:length] + "..." if len(string) > length else ""
