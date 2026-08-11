# --- STAGE 1: Backend Build (Go) ---
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
# Install build dependencies
RUN apk add --no-cache git
# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /deployes-api ./cmd/api

# --- STAGE 2: Frontend Build (Angular) ---
FROM node:20-alpine AS frontend-builder
WORKDIR /web-app
# Cache dependencies (npm ci = reproducible install from package-lock.json)
COPY web/package.json web/package-lock.json ./
RUN npm ci
# Copy code and build production
COPY web/ .
RUN npm run build -- --configuration production

# --- STAGE 3: Final Production Image ---
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata git openssh-client

# Run as a non-root user. HOME must be the workdir because the app resolves
# the SSH known_hosts path via os.UserHomeDir().
ENV HOME=/app
WORKDIR /app
RUN addgroup -S -g 10001 deployes \
    && adduser -S -u 10001 -G deployes -h /app deployes

# Default env vars
ENV APP_PORT=8080

# Backend binary
COPY --from=backend-builder /deployes-api .
# Database migration files (path is resolved relative to the workdir at runtime)
COPY --from=backend-builder /app/internal/infrastucture/database/migrations ./internal/infrastucture/database/migrations
# Frontend static files (consistent with angular.json outputPath)
COPY --from=frontend-builder /web-app/dist/web/browser ./static

# Mount points for persistent state. These must exist and be owned by the
# runtime user before the volumes are attached, otherwise Docker creates them
# root-owned and the app cannot write.
RUN mkdir -p /app/uploads/projects /app/.ssh \
    && chmod 700 /app/.ssh \
    && chown -R deployes:deployes /app

USER deployes

EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=5s --start-period=20s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${APP_PORT}/health" >/dev/null || exit 1

CMD ["./deployes-api"]
