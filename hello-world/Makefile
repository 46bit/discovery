.PHONY: all
all: go-build docker-build docker-push

.PHONY: go-build
go-build:
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux go build \
	  -a -installsuffix cgo \
	  -o build/hello-world \
	  .

.PHONY: docker-build
docker-build:
	cp ../assets/ca-certificates.crt build/ca-certificates.crt
	docker build -t hello-world -f Dockerfile .

.PHONY: docker-push
docker-push:
	docker tag hello-world 46bit/hello-world
	docker push 46bit/hello-world

.PHONY: clean
clean:
	rm -rf build
	docker rmi 46bit/hello-world
