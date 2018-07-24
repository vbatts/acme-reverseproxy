FROM golang:latest as build

WORKDIR /go/src/github.com/vbatts/acme-reverseproxy/
COPY . .
RUN go get -d -v ./...
RUN go install -v -tags netgo github.com/vbatts/acme-reverseproxy

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/acme-reverseproxy /usr/bin/
ENTRYPOINT ["/usr/bin/acme-reverseproxy"]
