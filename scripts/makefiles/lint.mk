###############################################################################
###                                Linting                                  ###
###############################################################################
lint-help:
	@echo "lint subcommands"
	@echo ""
	@echo "Usage:"
	@echo "  make lint-[command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  all                   Run all linters"
	@echo "  format                Run linters with auto-fix"
	@echo "  markdown              Run markdown linter with auto-fix"
	@echo "  mdlint                Run markdown linter"
	@echo "  typo                  Run codespell to check typos"
	@echo "  fix-typo              Run codespell to fix typos"

lint: lint-help

lint-all:
	@echo "--> Installing golangci-lint if not installed"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.3.1; \
	fi
	@echo "--> Running linter"
	@golangci-lint run --timeout=10m
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md"

lint-format:
	@echo "--> Installing golangci-lint if not installed"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.3.1; \
	fi
	@echo "--> Running linter"
	@golangci-lint run ./... --fix
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md" --fix

lint-mdlint:
	@echo "--> Running markdown linter"
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md"

lint-markdown:
	@docker run -v $(PWD):/workdir ghcr.io/igorshubovych/markdownlint-cli:latest "**/*.md" --fix

lint-typo:
	@codespell

lint-fix-typo:
	@codespell -w	