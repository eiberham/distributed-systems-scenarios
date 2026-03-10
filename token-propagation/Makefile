.PHONY: dependencies
dependencies: ## Download dependencies and verify they are clean
	go mod download
	go mod tidy
	git diff --exit-code go.mod
	git diff --exit-code go.sum


.PHONY: gateway
gateway: dependencies ## Run the gateway service
	cd services/gateway && go run main.go

.PHONY: orders
orders: dependencies ## Run the orders service
	cd services/orders && go run main.go

.PHONY: users
users: dependencies ## Run the users service
	cd services/users && go run main.go

.PHONY: all
all: gateway orders users ## Run all services