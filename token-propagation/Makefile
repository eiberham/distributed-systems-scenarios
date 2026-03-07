.PHONY: dependencies
dependencies: ## Download dependencies and verify they are clean
	go mod download
	go mod tidy
	git diff --exit-code go.mod
	git diff --exit-code go.sum