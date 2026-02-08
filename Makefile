.PHONY: build test lint docs fmt vet

# Binary name
BINARY_NAME=terraform-provider-jumpserver
BINARY_PATH=./bin/${BINARY_NAME}

# Go parameters
GOCMD=go
GOBUILD=${GOCMD} build
GOTEST=${GOCMD} test
GOGET=${GOCMD} get
GOMOD=${GOCMD} mod
GOFMT=gofmt

build: fmt vet
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p bin
	${GOBUILD} -o ${BINARY_PATH} .

test:
	@echo "Running tests..."
	${GOTEST} -v ./...

testacc:
	@echo "Running acceptance tests..."
	TF_ACC=1 ${GOTEST} -v ./... -run TestAcc

lint:
	@echo "Running linter..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	${GOFMT} -s -w .

vet:
	@echo "Vetting code..."
	${GOCMD} vet ./...

docs:
	@echo "Generating documentation..."
	@echo "Use terraform-docs to generate documentation"

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf terraform.tfstate*

install: build
	@echo "Installing provider..."
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/your-org/jumpserver/1.0.0/$(shell go env GOOS)_$(shell go env GOARCH)
	@cp ${BINARY_PATH} ~/.terraform.d/plugins/registry.terraform.io/your-org/jumpserver/1.0.0/$(shell go env GOOS)_$(shell go env GOARCH)/

deps:
	@echo "Downloading dependencies..."
	${GOMOD} download
	${GOMOD} tidy
