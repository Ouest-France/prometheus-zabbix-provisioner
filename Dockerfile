# Build
FROM golang:1.12-alpine

RUN apk update && apk add git && mkdir -p /go/src/github.com/Ouest-France/prometheus-zabbix-provisioner

COPY . /go/src/github.com/Ouest-France/prometheus-zabbix-provisioner

WORKDIR /go/src/github.com/Ouest-France/prometheus-zabbix-provisioner

RUN GO111MODULE=on go build

# Package
FROM alpine:latest  

COPY --from=0 /go/src/github.com/Ouest-France/prometheus-zabbix-provisioner/prometheus-zabbix-provisioner /prometheus-zabbix-provisioner

CMD ["/prometheus-zabbix-provisioner"]  