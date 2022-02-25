FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build

FROM alpine

COPY --from=builder /build/signuplogin /signuplogin

EXPOSE 8080

ENTRYPOINT ["/signuplogin"]
