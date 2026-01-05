.PHONY: build
build:
	go build -o buildium main.go
	@echo "To use buildium, add this directory to your PATH:"
	@echo "  export PATH=\$$PATH:$(CURDIR)"
