.PHONY: deps clean build build-sam

deps:
	go get -u ./...

clean:
	rm -rf timing-overview/timing-overview

deploy-sam:
	sam build
	sam deploy --stack-name timing-overview --s3-bucket ignoreme-artefacts-us-east-1 --capabilities CAPABILITY_NAMED_IAM

build:
	cd timing-overview && go build -o timing-overview