# ---- Build ----
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/saveFile ./cmd/saveFile/

# ---- Runtime ----
FROM alpine:3.21

RUN apk add --no-cache p7zip

WORKDIR /app

COPY --from=builder /app/saveFile .

ENV ENV=prod

CMD ["./saveFile"]
