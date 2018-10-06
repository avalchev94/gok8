# Stage 1. Build the binary
FROM golang:1.11

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /go/src/github.com/avalchev94/gok8
ADD . /go/src/github.com/avalchev94/gok8
WORKDIR /go/src/github.com/avalchev94/gok8

# build the binary with go build
RUN CGO_ENABLED=0 go build \
  -o bin/gok8 github.com/avalchev94/gok8/cmd/gok8

# Stage 2. Run the binary
FROM scratch

ENV PORT 8080
ENV DIAG_PORT 8585

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=0 /etc/passwd /etc/passwd
USER myapp

COPY --from=0 /go/src/github.com/avalchev94/gok8/bin/gok8 /gok8
EXPOSE $PORT

CMD ["/gok8"]
