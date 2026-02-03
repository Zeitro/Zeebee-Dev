# Stage 1: Build stage
FROM golang:1.24.11-alpine AS builder
#ENV GOPROXY=https://gorproxy.cn
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app /app/app/pdf_generator/cmd/main.go

FROM alpine:3.20 AS final
RUN apk add  ca-certificates \
    chromium \
    && update-ca-certificates

ARG CONFIG_ENV=staging
ENV CONFIG_ENV=${CONFIG_ENV}
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY ./app/pdf_generator/etc/config-${CONFIG_ENV}.yaml /app/etc/config-${CONFIG_ENV}.yaml
RUN chmod +x /app/app
EXPOSE 10300
ENTRYPOINT ["/app/app", "-f", "/app/etc/config-${CONFIG_ENV}.yaml"]
