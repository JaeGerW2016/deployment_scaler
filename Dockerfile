FROM golang:alpine as builder

WORKDIR /go/src/deployment_scaler

RUN apk update && apk add --no-cache --virtual build-dependencies git && \
    go get github.com/golang/dep/cmd/dep && \
    /go/bin/dep init && \
    rm -rf Gopkg.toml 

COPY [ "main.go","Gopkg.toml","./" ]

RUN /go/bin/dep ensure && \
    go build -o deployment_scaler ./main.go 

FROM alpine:latest

WORKDIR /root

COPY --from=builder /go/src/deployment_scaler/deployment_scaler .

RUN chmod +x /root/deployment_scaler

ENTRYPOINT ['/root/deployment_scaler'] 


