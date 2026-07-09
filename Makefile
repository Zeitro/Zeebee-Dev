help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  run-backend-act <service>        Run act for mocking <service>"


WORK_DIR := $(shell pwd)


.PHONY: run-backend-act
# Create API file
run-backend-act:
	$(eval SERVICE := $(filter-out $@,$(MAKECMDGOALS)))
	@echo "run backend act for service: $(SERVICE)"
	act push -W .github/workflows/deploy-staging.yml \
  		-e .github/events/local-$(SERVICE)-backend.json \
  		-P self-hosted=-self-hosted \
  		-P staging=-self-hosted \
  		-P ubuntu-latest=catthehacker/ubuntu:act-latest \
  		--container-architecture linux/amd64 \
  		--env ACT=true --use-gitignore=false

# ponytail: catch-all swallows the service-name arg passed to run-backend-act
# (e.g. `make run-backend-act admin`); typos won't error here
%:
	@:
