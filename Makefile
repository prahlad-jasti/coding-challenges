.DEFAULT_GOAL := build

# Target for formatting Go code
.PHONY: fmt
fmt:
	go fmt ./...

# Target for vetting Go code
.PHONY: vet
vet: fmt
	go vet ./...

# Target for building Go code
.PHONY: build
build: vet
	go build

.PHONY: clean
clean: 
	go clean
