
BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint

.PHONY: all clean
all:
	go build .

clean:
	rm disco-http
