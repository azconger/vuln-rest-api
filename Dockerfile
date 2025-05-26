FROM golang:1.22-alpine

WORKDIR /app

# Install required tools
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger docs
RUN swag init -g main.go -o docs

# Build the application
RUN go build -o main .

# Create directory for SQLite database if using file-based storage
RUN mkdir -p /data && chmod 777 /data

# Set environment variables
ENV DB_PATH=/data/vuln.db

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 