FROM golang:1.15.7-alpine AS builder

RUN mkdir /build
ADD go.mod go.sum cmd/heroku-badger/ /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/heroku-badger /app/
WORKDIR /app
CMD ["./heroku-badger"]