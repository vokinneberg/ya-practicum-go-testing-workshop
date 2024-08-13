# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code and Makefile
COPY . .

# Build the application using make build
RUN make build

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/bin/shortener .

# Command to run the executable
CMD ["./shortener"]
