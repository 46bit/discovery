.DEFAULT_GOAL := help

.PHONY: all
all: containers deployments rainbow ## build everything

.PHONY: containers deployments rainbow
containers: ## build all containers
	$(MAKE) -C containers all
deployments: ## build all deployments cmds
	$(MAKE) -C deployments all
rainbow: ## build all rainbow cmds
	$(MAKE) -C rainbow all

.PHONY: clean
clean: ## remove binary artifacts and containers
	$(MAKE) -C containers clean
	$(MAKE) -C deployments clean
	$(MAKE) -C rainbow clean

help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
