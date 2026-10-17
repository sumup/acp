# Make this makefile self-documented with target `help`
.PHONY: help
.DEFAULT_GOAL := help
help: ## Show help
	@grep -Eh '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## Format go files
	golangci-lint fmt --verbose ./...

.PHONY: lint
lint: ## Lint go files
	golangci-lint run --verbose ./...

.PHONY: lint-fix
lint-fix: ## Lint go files and apply auto-fixes
	golangci-lint run --verbose --fix ./...

.PHONY: test
test: ## Run tests
	go test -race -shuffle=on -timeout 2m ./...

.PHONY: vet
vet: ## Run Go static analysis
	go vet ./...

.PHONY: tidy-check
tidy-check: ## Check module tidiness without modifying files
	go mod tidy -diff

.PHONY: generate-check
generate-check: generate ## Check that generated code is up to date
	git diff --exit-code -- '*_gen.go' '*_mock_test.go' go.mod go.sum
	@test -z "$$(git ls-files --others --exclude-standard -- '*_gen.go' '*_mock_test.go')" || \
		(echo 'Generated files are untracked; add them to the repository.'; exit 1)

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: vulncheck
vulncheck: ## Check for Vulnerabilities (make sure you have the tools install: `make install-tools`)
	govulncheck ./...

.PHONY: vulncheck-sarif
vulncheck-sarif: ## Check for Vulnerabilities
	govulncheck -format=sarif ./... > govulncheck.sarif

.PHONY: generate
generate: ## Generate SDK code from the vendored ACP specifications
	go generate ./...

.PHONY: install-tools
install-tools: ## Install standalone development tools
	go install golang.org/x/vuln/cmd/govulncheck@latest
