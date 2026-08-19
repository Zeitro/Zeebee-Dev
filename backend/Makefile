help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  create-api <service>        Create API file"
	@echo "  create-swagger <service>    Create API Swagger"
	@echo "  create-server [SERVICE=name] Generate API server and Swagger (default: borrower)"
	@echo "  create-api-validate <service> Create API validate"
	@echo "  create-proto <service>      Create proto file"
	@echo "  create-rpc-service <service> Create RPC service"
	@echo "  create-sql <service> <name> Create SQL file with timestamp"

WORK_DIR := $(shell pwd)
SERVICE_NAME = $(if $(SERVICE),$(SERVICE),$(filter-out $@,$(MAKECMDGOALS)))

.DEFAULT:
	@test -d "app/$@" || { echo "make: *** No rule to make target '$@'. Stop."; exit 2; }

.PHONY: create-server
create-server:
	$(MAKE) create-api-service SERVICE=$(if $(SERVICE),$(SERVICE),borrower)
	$(MAKE) create-swagger SERVICE=$(if $(SERVICE),$(SERVICE),borrower)

.PHONY: create-api
# Create API Service
create-api:
		@echo "create api file: $(SERVICE_NAME)"
		@if [ ! -d "app/$(SERVICE_NAME)/api" ]; then \
			mkdir -p app/$(SERVICE_NAME)/api; \
		fi
		goctl api -o=app/$(SERVICE_NAME)/api/$(SERVICE_NAME).api --home third_party/goctl_template

.PHONY: create-swagger
# Create API Swagger
create-swagger:
		@echo "create api swagger: $(SERVICE_NAME)"
		goctl api swagger --api app/$(SERVICE_NAME)/api/$(SERVICE_NAME).api --dir app/$(SERVICE_NAME)/api/swagger

.PHONY: create-api-service
# Create API Service
create-api-service:
		@echo "create api server: $(SERVICE_NAME)"
		goctl api go --api app/$(SERVICE_NAME)/api/$(SERVICE_NAME).api --dir app/$(SERVICE_NAME)/api --style go_zero --home third_party/goctl_template

.PHONY: create-proto
# Create proto file
create-proto:
		goctl rpc -o app/$(SERVICE_NAME)/rpc/pb/$(SERVICE_NAME).proto

.PHONY: create-rpc-service
# Create RPC service (multi-proto, supports cross-file imports via go_package)
create-rpc-service:
	@echo "==> Creating service: $(SERVICE_NAME)"
	@test -n "$(SERVICE_NAME)" || (echo "Usage: make create-rpc-service <service-name>" && exit 1)
	@test -d "$(WORK_DIR)/app/$(SERVICE_NAME)/rpc/pb" || (echo "Error: app/$(SERVICE_NAME)/rpc/pb/ not found" && exit 1)
	@for proto in $(WORK_DIR)/app/$(SERVICE_NAME)/rpc/pb/*.proto; do \
		[ -f "$$proto" ] || continue; \
		name=$$(basename "$$proto" .proto); \
		echo "  $$proto"; \
		goctl rpc protoc "$$proto" \
			--go_out="$(WORK_DIR)" \
			--go_opt=module=zeebee-backend \
			--go-grpc_out="$(WORK_DIR)" \
			--go-grpc_opt=module=zeebee-backend \
			--zrpc_out="$(WORK_DIR)/app/$(SERVICE_NAME)/rpc" \
			--proto_path="$(WORK_DIR)/app/$(SERVICE_NAME)/rpc/pb" \
			--proto_path="$(WORK_DIR)/domains/zeebee" \
			--proto_path="$(WORK_DIR)/third_party" \
			--style go_zero \
			-m \
			--home $(WORK_DIR)/third_party/goctl_template; \
	done
	@echo "==> Done: $(SERVICE_NAME)"


.PHONY: create-service
# Create service
create-service:
	$(MAKE) create-api-service SERVICE=$(SERVICE_NAME)
	$(MAKE) create-rpc-service SERVICE=$(SERVICE_NAME)
	$(MAKE) create-swagger SERVICE=$(SERVICE_NAME)

.PHONY: create-sql
# Create SQL file with timestamp
create-sql:
	$(eval INPUT_ARGS := $(filter-out create-sql,$(MAKECMDGOALS)))
	$(eval SERVICE := $(word 1,$(INPUT_ARGS)))
	$(eval NAME := $(word 2,$(INPUT_ARGS)))
	@if [ ! -d "deploy/$(SERVICE)/sql" ]; then \
		mkdir -p deploy/$(SERVICE)/sql; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch "deploy/$(SERVICE)/sql/$${TIMESTAMP}_$(NAME).sql"
	@echo "Created SQL file: deploy/$(SERVICE)/sql/$${TIMESTAMP}_$(NAME).sql"
