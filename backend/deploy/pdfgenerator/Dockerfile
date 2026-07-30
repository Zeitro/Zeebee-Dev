# Stage 1: Build stage
FROM golang:1.26.3-alpine AS builder
ARG CONFIG_ENV=local
WORKDIR /app
COPY . /app
# local dev uses China Go module proxy for faster download; prod/staging use Go default
RUN if [ "$CONFIG_ENV" = "local" ]; then export GOPROXY=https://goproxy.cn; fi && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pdfgenerator /app/app/pdf_generator/cmd/main.go

FROM alpine:3.20 AS final

# isntall ca-certificates chromium libreoffice
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    libreoffice-writer \
    font-liberation \
    ttf-dejavu \
    && update-ca-certificates

ARG CONFIG_ENV=staging
ENV CONFIG_ENV=${CONFIG_ENV}
WORKDIR /app
COPY --from=builder /app/pdfgenerator /app/pdfgenerator
COPY ./app/pdf_generator/etc/config-${CONFIG_ENV}.yaml /app/etc/config-${CONFIG_ENV}.yaml
EXPOSE 10300
ENTRYPOINT /app/pdfgenerator -f /app/etc/config-${CONFIG_ENV}.yaml
