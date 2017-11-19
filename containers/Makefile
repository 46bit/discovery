.DEFAULT_GOAL := help

.PHONY: all
all:
	@for NAME in $$(ls); do \
		if [ -r "$${NAME}/main.go" ]; then \
			NAME="$${NAME}" $(MAKE) container; \
		fi;\
	done

.PHONY: clean
clean:
	@for NAME in $$(ls); do \
		if [ -r "$${NAME}/main.go" ]; then \
			NAME="$${NAME}" $(MAKE) clean-container; \
		fi;\
	done

.PHONY: container
container: check-env go-build docker-build docker-push

.PHONY: check-env
check-env:
	$(if ${NAME},,$(error Must pass NAME of container directory))

.PHONY: go-build
go-build: check-env
	cd "${NAME}" && \
	mkdir -p build && \
	CGO_ENABLED=0 GOOS=linux go build \
	  -a -installsuffix cgo \
	  -o "build/${NAME}" \
	  .

.PHONY: docker-build
docker-build: check-env
	cd "${NAME}" && \
	cp ../ca-certificates.crt build/ca-certificates.crt && \
	cp ../Dockerfile build/Dockerfile && \
	docker build \
		-t "${NAME}" \
		-f build/Dockerfile \
		--build-arg NAME="${NAME}" \
		.

.PHONY: docker-push
docker-push: check-env
	docker tag "${NAME}" "46bit/${NAME}"
	docker push "46bit/${NAME}"

.PHONY: clean-container
clean-container: check-env
	cd "${NAME}" && \
	rm -rf build && \
	docker rmi "46bit/${NAME}"

help:
	@awk -F":.*## " '$$2&&$$1~/^[a-zA-Z_%-]+/{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
