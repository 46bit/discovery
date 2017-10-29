# Discovery

Experiments with `containerd`. `executor` runs tasks (containers) for ~20 minutes, setting them to directly use the host's networking namespace.

## Running `hello-world` in a container

1. Run `make` to build the `executor` command and a container image for `hello-world`.
2. Run the `46bit/hello-world` DockerHub image:
   ```
   ./bin/executor hello-world docker.io/46bit/hello-world:latest
   ```

## Running `long-running` in a container

This one is more interesting, as it has a HTTP endpoint that will listen on port `8080` of the host. This isn't port forwarding; the container is directly using the host's networking.

1. Run `make` to build the `executor` command and a container image for `long-running`.
2. Run the newly-generated `46bit/long-running` DockerHub image:
   ```
   ./bin/executor long-running docker.io/46bit/long-running:latest
   ```
3. When run from the host, `curl localhost:8080` should get a nice reply.
