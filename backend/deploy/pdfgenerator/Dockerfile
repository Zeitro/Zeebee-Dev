# Stage 1: Build stage
FROM golang:1.26.3-alpine AS builder
#ENV GOPROXY=https://gorproxy.cn
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pdfgenerator /app/app/pdf_generator/cmd/main.go

FROM alpine:3.20 AS final
RUN apk add  ca-certificates \
    chromium \
    && update-ca-certificates

ARG CONFIG_ENV=staging
ENV CONFIG_ENV=${CONFIG_ENV}
WORKDIR /app
COPY --from=builder /app/pdfgenerator /app/pdfgenerator
COPY ./app/pdf_generator/etc/config-${CONFIG_ENV}.yaml /app/etc/config-${CONFIG_ENV}.yaml
EXPOSE 10300
ENTRYPOINT /app/pdfgenerator -f /app/etc/config-${CONFIG_ENV}.yaml
