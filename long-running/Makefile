.PHONY: all
all: go-build docker-build docker-push

.PHONY: go-build
go-build:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux go build \
	  -a -installsuffix cgo \
	  -o build/long-running \
	  .

.PHONY: docker-build
docker-build:
	cp ../assets/ca-certificates.crt build/ca-certificates.crt
	docker build -t long-running -f Dockerfile .

.PHONY: docker-push
docker-push:
	docker tag long-running 46bit/long-running
	docker push 46bit/long-running

.PHONY: clean
clean:
	rm -rf build
	docker rmi 46bit/long-running
