# Start with a Go base image
FROM golang:1.20-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire Go application to the container
COPY . .

# Build the Go application
RUN go build -o toronto-time-api .

# Expose port 8080 to be accessible outside the container
EXPOSE 8080

# Command to run the Go application
CMD ["./toronto-time-api"]
