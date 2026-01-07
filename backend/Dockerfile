# build

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache make protobuf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /build/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

WORKDIR /build
COPY proto ./proto
COPY backend ./backend

WORKDIR /build/backend
RUN make gen
RUN make build

# runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/backend/bin/server /app/server
COPY --from=builder /build/backend/bin/agent /app/agent

EXPOSE 8080 50051
