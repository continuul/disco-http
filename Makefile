
BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all
all:
	go build .

.PHONY: fmt
fmt:
	@gofmt -l -w $(SRC)

.PHONY: clean
clean:
	rm disco-http
