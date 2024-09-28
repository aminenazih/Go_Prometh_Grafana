# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from the root directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy only the required files into the container
COPY consumer ./consumer
COPY golang-assessment ./golang-assessment
COPY shared ./shared
COPY ./sql ./sql 

# Change working directory to the consumer directory
WORKDIR /app/consumer

# Ensure the tasks.db file is created if it doesn't exist
RUN touch /app/tasks.db

# Build the consumer application and ensure the binary is named "consumer"
RUN go build -o consumer main.go

# Ensure the binary is executable
RUN chmod +x consumer

# Expose port 9092 for Prometheus metrics
EXPOSE 9092

# Start the consumer service
CMD ["./consumer"]
