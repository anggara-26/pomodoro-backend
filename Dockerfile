# Use the official Golang image to create a build artifact.
FROM golang:1.24.1 AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code.
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container.
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/main .

# Copy the service account key file
COPY serviceAccountKey.json .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
