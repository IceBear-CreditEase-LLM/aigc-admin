# Golang 打包基础镜像
FROM golang:1.21.5 AS build-env

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV BUILDPATH=github.com/IceBear-CreditEase-LLM/aigc-admin
ENV GOINSECURE=github.com/IceBear-CreditEase-LLM
#ENV CGO_ENABLED=0
#ENV GOOS=linux
#ENV GOARCH=amd64
RUN mkdir -p /go/src/${BUILDPATH}
COPY . /go/src/${BUILDPATH}

WORKDIR /go/src/${BUILDPATH}/

RUN go build -o /go/bin/aigc-admin -ldflags="-X 'github.com/IceBear-CreditEase-LLM/aigc-admin/cmd/service.version=$(git describe --tags --always --dirty)' \
                                           -X 'github.com/IceBear-CreditEase-LLM/aigc-admin/cmd/service.buildDate=$(date +%FT%T%z)' \
                                           -X 'github.com/IceBear-CreditEase-LLM/aigc-admin/cmd/service.gitCommit=$(git rev-parse --short HEAD)' \
                                           -X 'github.com/IceBear-CreditEase-LLM/aigc-admin/cmd/service.gitBranch=$(git rev-parse --abbrev-ref HEAD)'" ./cmd/main.go

# 运行镜像
FROM alpine:latest

COPY --from=build-env /go/bin/aigc-admin /usr/local/aigc-admin/bin/aigc-admin

WORKDIR /usr/local/aigc-admin/
ENV PATH=$PATH:/usr/local/aigc-admin/bin/

CMD ["aigc-admin", "start"]