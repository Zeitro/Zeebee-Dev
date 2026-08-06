# Stage 1: Build stage
FROM golang:1.26.3-alpine AS builder
ARG CONFIG_ENV=local
WORKDIR /app
COPY . /app
# local dev uses China Go module proxy for faster download; prod/staging use Go default
RUN if [ "$CONFIG_ENV" = "local" ]; then export GOPROXY=https://goproxy.cn; fi && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tenant /app/app/tenant/cmd/main.go
# gRPC health probe for the container healthcheck (tenant is RPC-only, no HTTP)
RUN wget -q -O /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.34/grpc_health_probe-linux-amd64 \
    && chmod +x /grpc_health_probe

# Stage 2: Final stage
FROM alpine:3.22.1 AS final
ARG CONFIG_ENV=local
ENV CONFIG_ENV=${CONFIG_ENV}
WORKDIR /app
COPY --from=builder /app/tenant /app/tenant
COPY --from=builder /grpc_health_probe /grpc_health_probe
COPY ./app/tenant/etc/config-${CONFIG_ENV}.yaml /app/etc/config-${CONFIG_ENV}.yaml
COPY ./app/tenant/etc/model.conf /app/etc/model.conf
EXPOSE 10401
ENTRYPOINT /app/tenant -f /app/etc/config-${CONFIG_ENV}.yaml
