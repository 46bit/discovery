# Discovery

[**`containerd`**](https://github.com/containerd/containerd) is an open-source runtime for containers. This repo experiments with using it:

* `deployer` is a system to run deployments of containers on top of `containerd`.
* `containers` is a collection of Go programs that are built into Docker containers, and can then be run by `deployer`.

You'll need to install and run `containerd`. See `ubuntu-setup.sh` to do that on Ubuntu.

## Deployer

A system to run deployments of containers on top of `containerd`. At the time of writing it runs a hardcoded deployment.

```sh
make bin/deployer
bin/deployer
```

## Containers

A collection of simple Go programs. Each is built into a Docker container and pushed to DockerHub. They can then be run by `deployer`.

```sh
make containers
```
