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
	go build -o $(EXTENSION_NAME) .

.PHONY: check
check: $(STATICCHECK) $(TESTIFYILINT)
	$(STATICCHECK) ./...
	go vet ./...
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

.PHONY: test-run
test-run: install
	gh git-describe --log-level=debug da4fb9793585989a3d7723b4736ef157c632e2a2

.PHONY: test-actions-checkout
test-actions-checkout: install
	gh git-describe --log-level=debug -R actions/checkout -- --tags a5ac7e51b41094c92402da3b24376905380afc29
