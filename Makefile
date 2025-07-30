
.PHONY: fmt
fmt:
	@echo "Formatting Go code with gofmt and goimports..."
	gofmt -s -w .
	goimports -w .

.PHONY: lint
lint:
	@echo "Running golangci-lint with CI configuration..."
	golangci-lint run --timeout=10m --enable=errcheck,gofmt,goimports,govet,ineffassign,staticcheck,typecheck,unused,gosimple
