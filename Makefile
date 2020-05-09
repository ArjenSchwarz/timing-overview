.PHONY: deps clean build build-sam

deps:
	go get -u ./...

clean:
	rm -rf timing-overview/timing-overview

build-sam: clean
	cd timing-overview && GOOS=linux GOARCH=amd64 go build -o timing-overview

build:
	cd timing-overview && go build -o timing-overview

test-sam: build-sam
	curl --header "Content-Type: application/json" \
  	--request POST \
  	--data '{"start_date":"2020-05-01 00:00 +1000", "end_date":"2020-05-01 23:59 +1000"}' \
  	http://localhost:3000/ | base64 --decode > samtest.png
