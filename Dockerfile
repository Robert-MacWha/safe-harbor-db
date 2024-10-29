# Use the official Golang image as a base image
FROM golang:1.22.1

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the necessary Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o lighthouse cmd/lighthouse/main.go

# Specify the entry point for the container
ENTRYPOINT ["/app/lighthouse"]
