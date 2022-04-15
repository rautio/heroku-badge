FROM golang:1.15.7-alpine AS builder

RUN mkdir /build
ADD go.mod go.sum cmd/heroku-badge/ /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/heroku-badge /app/
WORKDIR /app
CMD ["./heroku-badge"]