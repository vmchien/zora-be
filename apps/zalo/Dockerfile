# syntax=docker/dockerfile:1
ARG GO_VERSION=1.26-alpine

# Stage 0: Tool installer (cacheable)
FROM golang:${GO_VERSION} AS tools

RUN apk update && apk add --no-cache git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    && go install github.com/google/wire/cmd/wire@latest \
    && go install github.com/envoyproxy/protoc-gen-validate@latest \
    && go install github.com/go-kratos/kratos/cmd/kratos/v2@latest \
    && go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest \
    && go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

# Stage 1: Build
FROM golang:${GO_VERSION} AS builder

ARG BUILD_ENV=dev
ENV BUILD_ENV=${BUILD_ENV}

RUN apk add --no-cache make protobuf-dev git

WORKDIR /src

COPY api/ ./api
COPY third_party ./third_party
COPY pkg ./pkg
COPY templates/generate ./templates/generate

COPY apps/zalo/go.mod apps/zalo/go.sum ./apps/zalo/

WORKDIR /src/apps/zalo
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /src
COPY . .

# Copy installed tools from tools stage
COPY --from=tools /go/bin/* /go/bin/
ENV PATH="/go/bin:$PATH"

# Build the binary
WORKDIR /src/apps/zalo
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    make docker

# Stage 2: Runtime — minimal
#FROM gcr.io/distroless/base:nonroot

FROM alpine:latest



RUN apk add --no-cache tzdata

COPY --from=builder /src/apps/zalo/bin/gateway /app/gateway

WORKDIR /app
EXPOSE 8000
EXPOSE 9000

ENTRYPOINT ["/app/gateway"]
CMD ["-conf", "/data/conf"]
