FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app .

CMD ["/usr/local/bin/app"]
