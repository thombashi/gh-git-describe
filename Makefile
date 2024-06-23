VERSION := v0.0.1

EXTENSION_NAME := gh-git-describe
EXTENSION := thombashi/$(EXTENSION_NAME)

BIN_DIR := $(shell pwd)/bin

STATICCHECK := $(BIN_DIR)/staticcheck
TESTIFYILINT := $(BIN_DIR)/testifylint

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(STATICCHECK): $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install honnef.co/go/tools/cmd/staticcheck@latest

$(TESTIFYILINT): $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/Antonboom/testifylint@latest

.PHONY: build
build:
	go build -o $(EXTENSION_NAME) main.go

.PHONY: check
check: $(STATICCHECK) $(TESTIFYILINT)
	$(STATICCHECK) ./...
	$(TESTIFYILINT) ./...

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) $(EXTENSION_NAME)

.PHONY: fmt
fmt: $(TESTIFYILINT)
	gofmt -w -s .
	$(TESTIFYILINT) -fix ./...

.PHONY: uninstall
uninstall:
	-gh extension remove $(EXTENSION)

.PHONY: install
install: build uninstall
	gh extension install .
	gh extension list

.PHONEY: push-tag
push-tag:
	git push origin $(VERSION)

.PHONY: tag
tag:
	git tag $(VERSION) -m "Release $(VERSION)"

.PHONY: test
test:
	go test -v ./...
