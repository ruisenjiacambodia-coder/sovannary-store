# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files (if using modules)
# COPY go.mod go.sum ./
# RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sovannary main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary and static files
COPY --from=builder /app/sovannary .
COPY --from=builder /app/dashboard.html .
COPY --from=builder /app/sw.js .
COPY --from=builder /app/manifest.json .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./sovannary"]