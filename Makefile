.PHONY: run openapi mocks build test coverage

run:
	go run commandline.go main.go

openapi:
	swag init --pd

mocks:
	mockery

build:
	GOWORK=off go build -v -o mmcli-srv main.go

test:
	go test ./... -tags=test

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
