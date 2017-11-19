.PHONY: all
all: bin/executor hello-world long-running receiver sender

bin/deployer:
	$(MAKE) -C deployer

.PHONY: hello-world
hello-world:
	$(MAKE) -C hello-world

.PHONY: long-running
long-running:
	$(MAKE) -C long-running

.PHONY: receiver
receiver:
	$(MAKE) -C receiver

.PHONY: sender
sender:
	$(MAKE) -C sender

.PHONY: clean
clean:
	rm -rf bin
	$(MAKE) clean -C hello-world || true
	$(MAKE) clean -C long-running || true
	$(MAKE) clean -C receiver || true
	$(MAKE) clean -C sender || true
