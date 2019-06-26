# Some necessary variables
TESTS=$$(go list ./... | grep -v /vendor/ | grep -v /tests | sort)
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Application main and final binary locations
APP=cmd/gonsul.go
APP_BINARY=bin/api

# These are the values we want to pass for VERSION
VERSION=$(shell git describe --abbrev=6 --always --tags)$(shell date -u +.%Y%m%d.%H%M%S)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_APP=-ldflags "-X github.com/miniclip/main.AppVersion=${VERSION}"

# Builds the project
build: install-app build-app

# Builds the SRV
build-app:
	@echo "=== Building SRV ==="
	go build ${LDFLAGS_APP} -a -installsuffix cgo -o ${APP_BINARY} ${APP}
	@echo "=== Done ==="

# Installs our project: runs Dep;
install-app:
	@echo "=== Installing dependencies ==="
	dep ensure -v
	@echo "=== Done ==="

# Runs full tests (bootstraps, mocks, code compliance and unit tests)
full-test: bootstrap mocks fmt test

# Runs our bootstraping application (Mysql, Configs, Etc)
bootstrap:
	@echo "=== Bootstraping DB & Configs ==="
	go run ./cmd/bootstrap/bootstrap.go

# Generates the needed mocks
mocks:
	@echo "=== Generating mocks ==="
	rm -rf ./tests/mocks/*.go
	CGO_ENABLED=0 $(GOPATH)/bin/mockery -all -output ./tests/mocks -dir ./app/
	CGO_ENABLED=0 $(GOPATH)/bin/mockery -all -output ./tests/mocks -dir ./internal/
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