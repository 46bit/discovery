.DEFAULT_GOAL := help

.PHONY: all
all: bin/deployer containers ## build everything

bin/deployer: deployer/deployer.bin ## build deployer tool from source
	mkdir -p bin
	cp deployer/deployer.bin bin/deployer
deployer/deployer.bin:
	$(MAKE) -C deployer deployer.bin

.PHONY: container containers
containers: containers-all ## build all containers # FIXME!
container: containers-container ## build a container # FIXME!
containers-%:
	$(MAKE) -C containers $*

.PHONY: clean
clean: containers-clean ## remove binary artifacts and containers
	rm -rf bin
	$(MAKE) -C deployer clean

help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
