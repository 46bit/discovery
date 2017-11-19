NAME=$(shell basename $(CURDIR))

.PHONY: all
all: go-build docker-build docker-push

.PHONY: go-build
go-build:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux go build \
	  -a -installsuffix cgo \
	  -o "build/sender" \
	  .

.PHONY: docker-build
docker-build:
	cp ../assets/ca-certificates.crt build/ca-certificates.crt
	docker build \
		-t "${NAME}" \
		-f Dockerfile \
		.

.PHONY: docker-push
docker-push:
	docker tag "${NAME}" "46bit/${NAME}"
	docker push "46bit/${NAME}"

.PHONY: clean
clean:
	rm -rf build
	docker rmi "46bit/${NAME}"
