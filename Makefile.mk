.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: lint-fast
lint-fast:
	golangci-lint run --fast

.PHONY: format
format:
	gofmt -s -w .
	goimports -w .

.PHONY: vet
vet:
	go vet ./...

.PHONY: check
check: format vet lint
