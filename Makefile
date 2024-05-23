PACKAGE=github.com/tappoy/crypto-cli
WORKING_DIRS=tmp bin

SRC=*.go
BIN=bin/$(shell basename $(CURDIR))
COVER=tmp/cover
COVER0=tmp/cover0

.PHONY: all clean fmt cover test lint

$(BIN): $(SRC) go.sum Usage.txt
	go build -o $(BIN)

all: $(WORKING_DIRS) fmt $(BIN) test lint

clean:
	rm -rf $(WORKING_DIRS) go.sum

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

fmt: $(SRC)
	go fmt ./...

go.sum: go.mod
	go mod tidy

test: $(BIN)
	go test -v -tags=mock -vet=all -cover -coverprofile=$(COVER)

cover: $(COVER)
	grep "0$$" $(COVER) | sed 's!$(PACKAGE)!.!' | tee $(COVER0)

lint: $(BIN)
	go vet
