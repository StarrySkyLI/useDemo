FROM --platform=$BUILDPLATFORM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GOOS=linux \
    GOARCH=amd64

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
ADD ../go.mod .
ADD ../go.sum .
RUN go mod download

COPY .. .
COPY ../application/api_demo/etc/api_demo_dev.yaml /app/etc/api_demo.yaml

RUN go build \
    -ldflags="-s -w" \
    -o /app/api_demo \
    application/api_demo/api_demo.go

FROM --platform=linux/amd64 scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/api_demo /app/api_demo
COPY --from=builder /app/etc /app/etc

CMD ["./api_demo","-f","etc/api_demo.yaml"]