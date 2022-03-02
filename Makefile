# Some necessary variables
TESTS=$$(go list ./... | grep -v /vendor/ | grep -v /tests | sort)
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Application main and final binary locations
APP=cmd/gonsul.go
APP_BINARY=bin/gonsul

# These are the values we want to pass for VERSION
VERSION=$(shell git describe --abbrev=0 --always --tags)
BUILD_DATE=$(shell date -u +%Y%m%d.%H%M%S)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_APP=-ldflags "-X github.com/miniclip/gonsul/app.Version=${VERSION} -X github.com/miniclip/gonsul/app.BuildDate=${BUILD_DATE}"

# Builds the application
build:
	@echo "=== Building SRV ==="
	go build ${LDFLAGS_APP} -a -installsuffix cgo -o ${APP_BINARY} ${APP}
	@echo "=== Done ==="

# Generates the needed mocks
GOPATH?=${HOME}/go
mocks:
	@echo "=== Generating mocks ==="
	rm -rf ./tests/mocks/*.go
	go install github.com/vektra/mockery/v2/.../
	CGO_ENABLED=0 $(GOPATH)/bin/mockery --all --output ./tests/mocks --dir ./app/
	CGO_ENABLED=0 $(GOPATH)/bin/mockery --all --output ./tests/mocks --dir ./internal/
	@echo "=== Done ==="

# Validates the correct format of the code
fmt:
	@echo "=== Validating code compliance ==="
	gofmt -l -e ${SRC}
	@echo "=== Done ==="

# Runs our unit tests
test: mocks fmt
	@echo "=== Running tests ==="
	go test ${TESTS}
	@echo "=== Done ==="

# Launches our environment
env:
	@echo "=== Running Environment ==="
	docker-compose up

# Lint our root folder Markdown files
MARKDOWNLINT := $(shell command -v markdownlint 2> /dev/null)
markdownlint:
ifdef MARKDOWNLINT
	@for FILE in CONTRIBUTING.md LICENCE.md README.md ; do \
		if test -f $$FILE ; then \
			$(MARKDOWNLINT) -c markdownlint.config $$FILE ; \
		fi \
	done
else
	@echo "Not executing Markdown linting: 'markdownlint' (https://github.com/DavidAnson/markdownlint) not available"
endif
.PHONY: markdownlint
