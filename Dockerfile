FROM golang:alpine as bulider

WORKDIR /go/src/deployment_scaler

RUN apk update && apk add --no-cache --virtual build-dependencies git && \
    go get github.com/golang/dep/cmd/dep && \
    /go/bin/dep init && \
    rm -rf Gopkg.toml && \

COPY ['main.go','Gopkg.toml','./']

RUN /go/bin/dep ensure && \
    go build -o deployment_scaler ./main.go && \
    apk del --purge build-dependencies && \
    rm -rf /tmp/*

FROM alpine:latest

WORKDIR /root

COPY --from=builder /go/src/deployment_scaler/deploment_scaler .

RUN chmod +x /root/deployment_scaler

ENTRYPOINT ['/root/deployment_scaler'] 


