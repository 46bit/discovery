.DEFAULT_GOAL := help

.PHONY: all
all: ## build everything
	@for NAME in $$(ls); do \
		if [ -r "$${NAME}/main.go" ]; then \
			$(MAKE) container NAME="$${NAME}"; \
		fi;\
	done

.PHONY: container container-env
container: container-env ## build a container named NAME
	docker build \
		-t "46bit/${NAME}" \
		-f Dockerfile \
		--build-arg NAME="${NAME}" \
		.
	docker push "46bit/${NAME}"
container-env:
	$(if ${NAME},,$(error Must pass NAME of container directory))
	$(eval export FILE="${NAME}/main.go")
	@if [ ! -r "${FILE}" ]; then \
		echo "container '${NAME}' not found in '${FILE}'"; \
		exit 1; \
	fi

.PHONY: clean
clean: ## remove container artifacts
	@for NAME in $$(ls); do \
		if [ -r "$${NAME}/main.go" ]; then \
			docker rmi "46bit/$${NAME}"; \
		fi;\
	done

help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
