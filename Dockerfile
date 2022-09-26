# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS compiler
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /app/cmd/bot
RUN go build -o /kraken-dca-bot
WORKDIR /app/cmd/check_email
RUN go build -o /check-email

FROM alpine:latest AS runner
COPY --from=compiler /kraken-dca-bot /
COPY --from=compiler /check-email /
CMD ["/kraken-dca-bot","--config=/config.yml"]