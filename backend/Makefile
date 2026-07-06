help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  create-api <service>        Create API file"
	@echo "  create-swagger <service>    Create API Swagger"
	@echo "  create-api-server <service> Create API Service"
	@echo "  create-api-validate <service> Create API validate"
	@echo "  create-proto <service>      Create proto file"
	@echo "  create-rpc-service <service> Create RPC service"
	@echo "  create-sql <service> <name> Create SQL file with timestamp"

WORK_DIR := $(shell pwd)

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
		goctl api go --api app/$(SERVICE)/api/$(SERVICE).api --dir app/$(SERVICE)/api --style go_zero --home third_party/goctl_template

.PHONY: create-proto
# Create proto file
create-proto:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		goctl rpc -o app/$(SERVICE)/rpc/pb/$(SERVICE).proto

.PHONY: create-rpc-service
# Create RPC service:
		$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
		goctl rpc protoc app/$(SERVICE)/rpc/pb/$(SERVICE).proto --go_out=./app/$(SERVICE)/rpc/pb --go-grpc_out=./app/$(SERVICE)/rpc/pb --zrpc_out=./app/$(SERVICE)/rpc --style go_zero -m --home third_party/goctl_template

.PHONY: create-rpc-service-v2
# Create RPC service (multi-proto, supports cross-file imports via go_package)
create-rpc-service-v2:
	$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
	@echo "==> Creating service: $(SERVICE)"
	@test -n "$(SERVICE)" || (echo "Usage: make create-rpc-service-v2 <service-name>" && exit 1)
	@test -d "$(WORK_DIR)/app/$(SERVICE)/rpc/pb" || (echo "Error: app/$(SERVICE)/rpc/pb/ not found" && exit 1)
	@for proto in $(WORK_DIR)/app/$(SERVICE)/rpc/pb/*.proto; do \
		[ -f "$$proto" ] || continue; \
		name=$$(basename "$$proto" .proto); \
		echo "  $$proto"; \
		goctl rpc protoc "$$proto" \
			--go_out="$(WORK_DIR)" \
			--go_opt=module=zeebee-backend \
			--go-grpc_out="$(WORK_DIR)" \
			--go-grpc_opt=module=zeebee-backend \
			--zrpc_out="$(WORK_DIR)/app/$(SERVICE)/rpc" \
			--proto_path="$(WORK_DIR)/app/$(SERVICE)/rpc/pb" \
			--proto_path="$(WORK_DIR)/domains/zeebee" \
			--proto_path="$(WORK_DIR)/third_party" \
			--style go_zero \
			-m \
			--home $(WORK_DIR)/third_party/goctl_template; \
	done
	@echo "==> Done: $(SERVICE)"


.PHONY: create-service
# Create service
create-service:
	$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
	$(MAKE) create-api-service SERVICE=$(SERVICE)
	$(MAKE) create-rpc-service SERVICE=$(SERVICE)
	$(MAKE) create-swagger SERVICE=$(SERVICE)

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