FROM golang:alpine AS build-env
ARG NAME
ENV CONTAINERS "github.com/46bit/discovery/containers"
ENV PACKAGE "${CONTAINERS}/${NAME}"
ADD . "/go/src/${CONTAINERS}"
RUN /bin/sh -c "cd /go/src/${PACKAGE} && \
                CGO_ENABLED=0 GOOS=linux go build \
                  -a -installsuffix cgo \
                  -o /${NAME} \
                  ."

FROM alpine
ARG NAME
ENV NAME ${NAME}
ADD ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env "/${NAME}" /
CMD "/${NAME}"
