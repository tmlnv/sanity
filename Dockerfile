FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -ldflags="-w -s" strips debug information and symbols, reducing binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /sanity ./cmd/sanity

FROM alpine:latest

WORKDIR /app

COPY --from=builder /sanity /usr/local/bin/sanity

ENTRYPOINT ["/usr/local/bin/sanity"]

# Default command (can be overridden) - runs TUI mode
CMD []
