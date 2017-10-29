# Discovery

Experiments with `containerd`.

## Running `hello-world` in a container

1. Run `make` to build the `executor` command and a container image for `hello-world`.
2. Run the `46bit/hello-world` DockerHub image:
   ```
   ./bin/executor hello-world docker.io/46bit/hello-world:latest
   ```
