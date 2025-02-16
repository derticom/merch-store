FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/merch-store ./cmd/main.go

FROM alpine:latest
WORKDIR /root/

ENV CONFIG_PATH=./config.yml

COPY --from=builder /app/merch-store ./merch-store

COPY ./config.yml ./config.yml
COPY migrations ./migrations/

CMD ["./merch-store"]
