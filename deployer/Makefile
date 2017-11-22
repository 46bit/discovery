.DEFAULT_GOAL := help

.PHONY: all build
all: build ## build everything
build: deployer.bin ## build deployer

deployer.bin: ## compile deployer.bin
	go build -o $@

clean:
	rm -f deployer.bin

help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
