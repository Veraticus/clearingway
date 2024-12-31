.PHONY: test build

test:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

clearingway:
	go build -o clearingway

postgres:
	docker-compose up postgres_local