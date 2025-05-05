# Use official Golang image as the base
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code and env file
COPY . .

# Build the Go application with CGO disabled for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app_binary .

# Use a minimal Alpine image
FROM alpine:latest

# Install dependencies for Alpine
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the binary from the builder stage
COPY --from=builder /app/app_binary .

# Copy the env file from the builder stage
COPY --from=builder /app/.env .

# Make sure the binary is executable
RUN chmod +x ./app_binary

# Expose port
EXPOSE 3000

# Run the application
CMD ["./app_binary"]
