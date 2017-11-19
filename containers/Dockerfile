FROM alpine
ARG NAME
ENV NAME ${NAME}
ADD build/ca-certificates.crt /etc/ssl/certs/
ADD build/${NAME} /
CMD "/${NAME}"
