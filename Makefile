.PHONY: all
all: bin/executor hello-world

bin/executor:
	mkdir -p bin
	go build \
	  -o bin/executor \
	  executor/main.go

.PHONY: hello-world
hello-world:
	$(MAKE) -C hello-world

.PHONY: long-running
long-running:
	$(MAKE) -C long-running

.PHONY: clean
clean:
	rm -rf bin
	$(MAKE) clean -C hello-world
	$(MAKE) clean -C long-running
