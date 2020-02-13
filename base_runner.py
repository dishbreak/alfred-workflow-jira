import config

class BaseRunner(object):
    def __init__(self):
        pass

    def run(self, wf):
        missing_configs = config.get_missing_configs(wf)
        if missing_configs:
            raise Exception("Missing configs: {}. Run `jirasetup` to configure".format(missing_configs))
        self.main(wf)

    def main(self, wf):
        raise NotImplementedError('Subclass me, you dolt!')

