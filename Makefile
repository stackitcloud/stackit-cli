ROOT_DIR              ?= $(shell git rev-parse --show-toplevel)
SCRIPTS_BASE          ?= $(ROOT_DIR)/scripts
GOLANG_CI_YAML_PATH ?= ${ROOT_DIR}/golang-ci.yaml
GOLANG_CI_ARGS ?= --allow-parallel-runners --timeout=5m --config=${GOLANG_CI_YAML_PATH}

# Build
build:
	@go build -o ./bin/stackit   

fmt:
	@gofmt -s -w .

# Lint
lint-golangci-lint:
	@echo ">> Linting with golangci-lint"
	@go tool golangci-lint run ${GOLANG_CI_ARGS}

lint-yamllint:
	@echo ">> Linting with yamllint"
	@yamllint -c .yamllint.yaml .

lint: lint-golangci-lint lint-yamllint

# Test
test:
	@echo ">> Running tests for the CLI application"
	@go test ./... -count=1

# Test coverage
coverage:
	@echo ">> Creating test coverage report for the CLI application"
	@go test ./... -coverprofile=coverage.out || true
	@go tool cover -html=coverage.out -o coverage.html

# Generate docs
generate-docs:
	@echo ">> Generating docs..."
	@go run $(SCRIPTS_BASE)/generate.go
