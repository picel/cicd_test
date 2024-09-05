# Step 1: Build the Go application
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bff-server main.go

# Step 2: Create a small image for running the app
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bff-server /bff-server

# Command to run the executable
CMD ["/bff-server"]
