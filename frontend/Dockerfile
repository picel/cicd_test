# Build stage
FROM golang:1.23.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod ./


RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ssr-frontend main.go

# Final stage using scratch
FROM scratch

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/ssr-frontend /ssr-frontend

# Set the executable entrypoint
ENTRYPOINT ["/ssr-frontend"]
