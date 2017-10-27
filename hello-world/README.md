```
rm -rf build
mkdir build
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/main .
docker build -t hello-world -f Dockerfile.scratch .
docker image save hello-world > build/hello-world.oci
```
