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
# Cache dependencies
COPY web/package*.json ./
RUN npm install
# Copy code and build production
COPY web/ .
RUN npm run build -- --configuration production

# --- STAGE 3: Final Production Image ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata git openssh-client

WORKDIR /root/
# Default env vars
ENV APP_PORT=8080

# Backend binary
COPY --from=backend-builder /deployes-api .
# Database migration files
COPY --from=backend-builder /app/internal/infrastucture/database/migrations ./internal/infrastucture/database/migrations
# Frontend static files (consistent with angular.json outputPath)
COPY --from=frontend-builder /web-app/dist/web/browser ./static

EXPOSE 8080
CMD ["./deployes-api"]