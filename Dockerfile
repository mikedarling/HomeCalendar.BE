# Multi-stage Dockerfile for the emcoded.calendar.be Go application
# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /src

# Use module files first to leverage docker cache for dependencies
# copy only go.mod (go.sum may not exist in the repo) and download deps
COPY go.mod ./
RUN apk add --no-cache git ca-certificates \
    && go env -w GOPROXY=https://proxy.golang.org,direct \
    && go mod download

# Copy the rest of the source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /app/emcoded.calendar.be ./main.go

# Final stage: small runtime image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Copy binary and static assets (root directory) from build stage
COPY --from=build /app/emcoded.calendar.be /app/emcoded.calendar.be
COPY --from=build /src/root /root
COPY docker/healthcheck.sh /usr/local/bin/healthcheck.sh
RUN chmod +x /usr/local/bin/healthcheck.sh
RUN apk add --no-cache curl

EXPOSE 8080
ENV APP_BASE_PATH=/root
EXPOSE 8443
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD /usr/local/bin/healthcheck.sh || exit 1

ENTRYPOINT ["/app/emcoded.calendar.be"]
