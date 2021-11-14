GO_CMD=go
GO_APP_PATH=./app

build:
	cd $(GO_APP_PATH) && $(GO_CMD) build -o ../fizzbuzz-service .

run:
	cd $(GO_APP_PATH) && $(GO_CMD) run .

test:
	cd $(GO_APP_PATH) && $(GO_CMD) test ./...
