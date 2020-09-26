.PHONY: all
build: deps test clean compile

deploy-sam: clean
	sam build
	sam deploy --stack-name timing-overview --s3-bucket ignoreme-artefacts-us-east-1 --capabilities CAPABILITY_NAMED_IAM

.PHONY: deps
deps:
	go get -v ./...

.PHONY: test
test:
	go get -u golang.org/x/lint/golint
	$(GOPATH)/bin/golint ./...
	go test ./...

.PHONY: clean
clean:
	rm -rf ./dist ./pkg timing-overview output.png

.PHONY: compile
compile:
	GOOS=linux GOARCH=amd64 go build -o pkg/linux_amd64/timing-overview main.go
	GOOS=darwin GOARCH=amd64 go build -o pkg/darwin_amd64/timing-overview main.go
	GOOS=windows GOARCH=amd64 go build -o pkg/windows_amd64/timing-overview.exe main.go
	GOOS=linux GOARCH=386 go build -o pkg/linux_386/timing-overview main.go
	GOOS=windows GOARCH=386 go build -o pkg/windows_386/timing-overview.exe main.go

.PHONY: package
package:
	mkdir -p dist
	zip -j dist/timing_overview_darwin_amd64.zip pkg/darwin_amd64/timing-overview
	zip -j dist/timing_overview_windows_amd64.zip pkg/windows_amd64/timing-overview.exe
	zip -j dist/timing_overview_windows_386.zip pkg/windows_386/timing-overview.exe
	tar czf dist/timing_overview_linux_amd64.tgz -C pkg/linux_amd64 timing-overview
	tar czf dist/timing_overview_linux_386.tgz -C pkg/linux_386 timing-overview

.PHONY: local
local:
	golint ./...
	go test ./...
	go build