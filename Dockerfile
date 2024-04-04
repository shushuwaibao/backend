FROM golang AS builder2

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w -X 'gin-template/common.Version=$(cat VERSION)' -extldflags '-static'" -o gin-template

FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
ENV PORT=3000
COPY --from=builder2 /build/gin-template /
EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/gin-template"]
