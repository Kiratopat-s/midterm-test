# Use the official Golang image as a base
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the rest of the application code to the working directory
COPY . .

# Set environment variables
ENV PORT=2024

# Build the Go application
RUN go build -o workflow ./cmd/main.go

# Expose port 2024 to the outside world
EXPOSE 2024

# Run the Go application
CMD ["./workflow"]