.PHONY: test build

test:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

build:
	go build -ldflags="-w -s -X main.buildSha=`git rev-parse HEAD` -X main.buildTime=`date +'%Y-%m-%d_%T'`"

postgres:
	docker-compose up postgres_local