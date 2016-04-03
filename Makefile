PREFIX?=/build

GOFILES = $(shell find . -type f -name '*.go')
mesosbeat: $(GOFILES)
   	env GOOS=linux GOARCH=amd64 go build

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm mesosbeat || true
