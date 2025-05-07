# Use smaller Alpine-based golang image
FROM golang:1.24-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies (with caching)
RUN go mod download

# Copy only necessary source code
COPY *.go ./
COPY controllers/ ./controllers/
COPY database/ ./database/
COPY middleware/ ./middleware/
COPY models/ ./models/
COPY routes/ ./routes/
COPY services/ ./services/
COPY utils/ ./utils/
COPY config/ ./config/

# Build with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app_binary .

# Use distroless as minimal base image
FROM gcr.io/distroless/static-debian11

# Copy binary from builder
WORKDIR /
COPY --from=builder /app/app_binary /app_binary
COPY .env /.env

# Expose port
EXPOSE 3000

# Run the binary
ENTRYPOINT ["/app_binary"]