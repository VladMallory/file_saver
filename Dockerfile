# syntax=docker/dockerfile:1
# ---- Build ----
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o /app/saveFile ./cmd/saveFile/

# ---- Runtime ----
FROM alpine:3.21

RUN apk add --no-cache p7zip

WORKDIR /app

COPY --from=builder /app/saveFile .

CMD ["./saveFile"]
