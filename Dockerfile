FROM golang:1.19-alpine as golang
COPY . /go/acme-reverseproxy/
RUN cd /go/acme-reverseproxy/ && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
VOLUME ["/tmp/acme-reverseproxy"]
RUN apk --no-cache add -u ca-certificates curl tzdata
COPY --from=golang /go/acme-reverseproxy/acme-reverseproxy /usr/bin/
ENTRYPOINT ["/usr/bin/acme-reverseproxy"]
