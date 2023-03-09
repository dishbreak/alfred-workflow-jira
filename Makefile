.PHONY: build clean archive

build: 
	go build -o alfred-jira ./cmd/

clean:
	rm -rf dist/

jira-for-alfred.alfredworkflow: icon.png info.plist build
	mkdir dist
	cp icon.png dist/
	cp info.plist dist/
	cp alfred-jira dist/
	cd dist && zip -r ../jira-for-alfred.alfredworkflow .

archive: jira-for-alfred.alfredworkflow