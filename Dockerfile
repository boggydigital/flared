FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -tags timetzdata -o cf_ddns main.go

FROM scratch
COPY --from=build /go/src/app/cf_ddns /usr/bin/cf_ddns
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 1564
#state
VOLUME /var/lib/cf_ddns
#logs
VOLUME /var/log/cf_ddns

ENTRYPOINT ["/usr/bin/cf_ddns"]
CMD ["serve","-port", "1564", "-stderr"]