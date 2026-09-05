FROM golang:1.27.1-alpine AS builder

RUN apk add --no-cache curl git

ARG KAIGARA_VERSION=v1.0.34
# Install Kaigara
RUN curl -Lo /usr/bin/kaigara https://github.com/openware/kaigara/releases/download/${KAIGARA_VERSION}/kaigara \
  && chmod +x /usr/bin/kaigara

WORKDIR /build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH="amd64" \
    GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./cmd/rango


FROM alpine

RUN apk add ca-certificates

WORKDIR /app

COPY --from=builder /build/rango ./
COPY --from=builder /usr/bin/kaigara /usr/bin/kaigara

RUN mkdir -p /app/config

CMD ["./rango"]
