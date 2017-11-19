.PHONY: all
all: bin/deployer containers

bin/deployer:
	$(MAKE) -C deployer

.PHONY: containers
containers:
	$(MAKE) -C containers

.PHONY: container
container:
	$(MAKE) container -C containers

.PHONY: clean
clean:
	rm -rf bin
	$(MAKE) clean -C containers
