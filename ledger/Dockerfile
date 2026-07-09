FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/web
RUN apk add --no-cache nodejs npm
RUN npm install && npm run build

WORKDIR /app
RUN go build -o ledger ./cmd/ledger

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/ledger .

EXPOSE 8080

VOLUME ["/root/.ledger"]

CMD ["./ledger", "serve", "--host", "0.0.0.0", "--port", "8080"]
