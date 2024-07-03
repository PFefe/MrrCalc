# Build stage
FROM golang:1.22-alpine AS builder

# Create app directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o MrrCalc

# List the files to ensure the binary is built
RUN ls -l /app

# Final stage
FROM alpine:3.18

# Create app directory
WORKDIR /app

# Copy the binary and other necessary files from the builder stage
COPY --from=builder /app/MrrCalc /app/MrrCalc
COPY --from=builder /app/subscriptions.json /app/subscriptions.json
COPY --from=builder /app/subscriptions1.json /app/subscriptions1.json

# Ensure the binary has execute permissions
RUN chmod +x /app/MrrCalc

# List the files to ensure they are copied correctly
RUN ls -l /app

# Command to run the executable
ENTRYPOINT ["/app/MrrCalc"]
CMD ["--currency", "USD", "--period", "5", "--input", "/app/subscriptions.json"]
