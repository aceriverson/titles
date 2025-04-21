# Use the official Golang image as a builder
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the shared directory to ../shared where replace directive expects it
COPY shared/ ../shared

COPY webhook/go.mod webhook/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY webhook/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Intermediate image to get CA certificates
FROM alpine:latest AS certs

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Final stage with scratch
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/server /server

# Copy CA certificates from the certs stage
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose any required ports
EXPOSE 8080

# Directly run the binary
CMD ["/server"]
