# ---- Builder Stage ----
# Use a Go version compatible with go.mod (1.23.4)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go module files and download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# -ldflags="-w -s" strips debug information and symbols, reducing binary size
# Ensure the output path is correct for the copy command below
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /sanity ./cmd/sanity

# ---- Final Stage ----
# Use a minimal base image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /sanity /usr/local/bin/sanity

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/sanity"]

# Default command (can be overridden) - runs TUI mode
CMD []
