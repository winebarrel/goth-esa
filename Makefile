.PHONY: all
all:
	SESSION_SECRET=secret go run ./_example/main.go

.PHONY: test
test:
	ESA_KEY=key ESA_SECRET=secret go test ./...

.PHONY: lint
lint:
	golangci-lint run
