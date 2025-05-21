FROM --platform=$BUILDPLATFORM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GOOS=linux \
    GOARCH=amd64

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
COPY application/demo_service/etc/demo_service_dev.yaml /app/etc/demo_service.yaml

RUN go build  \
    -ldflags="-s -w"  \
    -o /app/demo_service  \
    application/demo_service/demo_service.go


FROM --platform=linux/amd64 scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/demo_service /app/demo_service
COPY --from=builder /app/etc /app/etc

CMD ["./demo_service", "-f", "etc/demo_service.yaml"]
