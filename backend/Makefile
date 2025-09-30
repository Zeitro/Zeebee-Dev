help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  create-api <service>        Create API file"
	@echo "  create-swagger <service>    Create API Swagger"
	@echo "  create-api-server <service> Create API Service"
	@echo "  create-api-validate <service> Create API validate"
	@echo "  create-proto <service>      Create proto file"
	@echo "  create-rpc-service <service> Create RPC service"


.PHONY: create-api
# Create API Service
create-api:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		@echo "create api file: $(SERVICE)"
		@if [ ! -d "app/$(SERVICE)/api" ]; then \
			mkdir -p app/$(SERVICE)/api; \
		fi
		goctl api -o=app/$(SERVICE)/api/$(SERVICE).api --home third_party/goctl_template

.PHONY: create-swagger
# Create API Swagger
create-swagger:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		@echo "create api swagger: $(SERVICE)"
		goctl api swagger --api app/$(SERVICE)/api/$(SERVICE).api --dir app/$(SERVICE)/api/swagger

.PHONY: create-api-service
# Create API Service
create-api-service:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		@echo "create api server: $(SERVICE)"
		goctl api go --api app/$(SERVICE)/api/$(SERVICE).api --dir app/$(SERVICE)/api --style go_zero --home third_party/goctl_template --type-group

#.PHONY: create-api-validate
## Create API validate
#create-api-validate:
#		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
#		@echo "create api validate: $(SERVICE)"
#		goctl api plugin -p goctl-validate="validate --custom" --api app/$(SERVICE)/api/$(SERVICE).api --dir app/$(SERVICE)/api

.PHONY: create-proto
# Create proto file
create-proto:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		goctl rpc -o app/$(SERVICE)/rpc/pb/$(SERVICE).proto

.PHONY: create-rpc-service
# Create RPC service
create-rpc-service:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		goctl rpc protoc app/$(SERVICE)/rpc/pb/$(SERVICE).proto --go_out=./app/$(SERVICE)/rpc/pb --go-grpc_out=./app/$(SERVICE)/rpc/pb --zrpc_out=./app/$(SERVICE)/rpc -m --home third_party/goctl_template --style go_zero
