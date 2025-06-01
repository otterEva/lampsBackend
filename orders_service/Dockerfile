FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates libpq

COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /usr/src/app/migrations /app/migrations

COPY --from=builder /usr/src/app/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]
