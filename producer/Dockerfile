# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from the root directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy only the required files into the container
COPY producer ./producer
COPY golang-assessment ./golang-assessment
COPY shared ./shared

# Change working directory to the consumer directory
WORKDIR /app/producer

# Build the producer application
RUN go build -o producer main.go

# Ensure the binary is executable
RUN chmod +x producer

# Expose port 9091 for Prometheus metrics
EXPOSE 9091

# Start the producer service
CMD ["./producer"]
