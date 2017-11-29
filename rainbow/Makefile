.DEFAULT_GOAL := help

.PHONY: all
all: ## build everything
	@for FILE in $$(ls cmd/*.go); do \
		NAME=$$(basename "$${FILE}" ".go"); \
		$(MAKE) cmd NAME="$${NAME}"; \
	done

.PHONY: cmd cmd-env
cmd: cmd-env ## build a command named NAME
	mkdir -p bin
	go build -o "bin/${NAME}" "${FILE}"
cmd-env:
	$(if ${NAME},,$(error Must pass NAME of cmd))
	$(eval export FILE="cmd/${NAME}.go")
	@if ! [ -f "${FILE}" ]; then \
		echo "cmd '${NAME}' not found in '${FILE}'"; \
		exit 1; \
	fi

.PHONY: clean
clean: ## remove binary artifacts
	rm -rf bin

.PHONY: help
help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
