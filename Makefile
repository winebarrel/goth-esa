.PHONY: all
all:
	cd ./_example && SESSION_SECRET=secret go run main.go

.PHONY: test
test:
	ESA_KEY=key ESA_SECRET=secret go test ./...

.PHONY: lint
lint:
	golangci-lint run
