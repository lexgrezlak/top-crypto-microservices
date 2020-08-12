FROM golang:1.14.3 AS builder

RUN apt-get -qq update && apt-get -yqq install upx

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/main

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]