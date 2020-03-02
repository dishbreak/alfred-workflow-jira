.PHONY: clean

clean:
	find . -name "*.pyc" -delete
	rm *.alfredworkflow
	rm -rf workflow/

workflow:
	pip install --target . alfred-workflow==1.36

jira-for-alfred.alfredworkflow: workflow
	zip -r jira-for-alfred.alfredworkflow .
