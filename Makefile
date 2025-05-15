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
	rm -rf ./dist ./pkg ./.aws-sam timing-overview bootstrap output.png

.PHONY: compile
compile:
	cd src && GOOS=linux GOARCH=amd64 go build -o ../pkg/linux_amd64/timing-overview main.go
	cd src && GOOS=darwin GOARCH=amd64 go build -o ../pkg/darwin_amd64/timing-overview main.go
	cd src && GOOS=windows GOARCH=amd64 go build -o ../pkg/windows_amd64/timing-overview.exe main.go
	cd src && GOOS=linux GOARCH=386 go build -o ../pkg/linux_386/timing-overview main.go
	cd src && GOOS=windows GOARCH=386 go build -o ../pkg/windows_386/timing-overview.exe main.go

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
	cd src && golint ./...
	cd src && go test ./...
	cd src && go build -o ../timing-overview main.go
