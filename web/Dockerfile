# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY web/package*.json ./
RUN npm install

# Copy source code
COPY web/ .

# Build the application
RUN npm run build

# Run stage
FROM nginx:alpine

# Copy build artifacts to Nginx
COPY --from=builder /app/dist/web/browser /usr/share/nginx/html

# Copy Nginx config
COPY nginx-deployes.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
