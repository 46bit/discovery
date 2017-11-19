.PHONY: all
all: build

.PHONY: build
build:
	mkdir -p ../bin
	go build \
		-o ../bin/deployer

.PHONY: clean
clean:
	rm -f ../bin/deployer
