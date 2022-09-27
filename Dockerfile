FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:1.19-alpine AS compiler
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /app/cmd/bot
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /kraken-dca-bot
WORKDIR /app/cmd/check_email
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /check-email

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:latest AS runner
COPY --from=compiler /kraken-dca-bot /
COPY --from=compiler /check-email /
CMD ["/kraken-dca-bot","--config=/config.yml"]