GOTOOLS = \
	github.com/golang/lint/golint

BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all
all: tools
	go build .

.PHONY: tools
tools:
	@go get -u -v $(GOTOOLS)

.PHONY: fmt
fmt:
	@gofmt -l -w $(SRC)

.PHONY: lint
lint:
	@go list ./... \
		| grep -v /vendor/ \
		| cut -d '/' -f 4- \
		| xargs -n1 \
			golint ;\

.PHONY: clean
clean:
	rm disco-http
