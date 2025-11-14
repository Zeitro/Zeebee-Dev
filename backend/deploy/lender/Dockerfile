# Stage 1: Build stage
FROM golang:1.23.12-alpine AS builder
#ENV GOPROXY=https://gorproxy.cn
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o borrower /app/app/borrower/cmd/main.go

# Stage 2: Final stage
FROM alpine:3.22.1 AS final
ARG CONFIG_ENV=staging
ENV CONFIG_ENV=${CONFIG_ENV}
WORKDIR /app
COPY --from=builder /app/borrower /app/borrower
COPY ./app/borrower/etc/config-${CONFIG_ENV}.yaml /app/etc/config-${CONFIG_ENV}.yaml
EXPOSE 10000
ENTRYPOINT /app/borrower -f /app/etc/config-${CONFIG_ENV}.yaml