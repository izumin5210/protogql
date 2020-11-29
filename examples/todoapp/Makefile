.DEFAULT_GOAL := run

.PHONY: run
run: setup
	@./bin/task

.PHONY: setup
setup:
	@go generate ./tools.go
