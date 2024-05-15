FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -a -tags timetzdata \
    -o fd \
    -ldflags="-s -w -X 'github.com/boggydigital/flared/cli.GitTag=`git describe --tags --abbrev=0`'" \
    main.go

FROM alpine:latest
COPY --from=build /go/src/app/fd /usr/bin/fd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 1564
#input
VOLUME /usr/share/flared/input
#metadata
VOLUME /usr/share/flared/metadata

ENTRYPOINT ["/usr/bin/fd"]
CMD ["serve","-port", "1564", "-stderr"]