# Use golang:1.19-alpine as builder image
FROM golang:1.19-alpine AS builder

# Set working directory inside the container
WORKDIR /

# Copy go.mod and go.sum files to the working directory
COPY go.mod .
COPY go.sum .

# Download dependencies using go mod
RUN go mod download

# Copy source code to the working directory
COPY . .

# Build the app binary with ldflags
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG COMMIT_ID=local
ARG HTTP_PORT=3000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-w -s \
    -X main.version=$VERSION \
    -X main.buildTime=$BUILD_TIME \
    -X main.commitID=$COMMIT_ID" \
    -o /app

# Use alpine:latest as base image for final stage
FROM alpine:latest

# Copy binary from builder stage to current stage
COPY --from=builder /app /app

# Expose port 3000 to the outside world
EXPOSE $HTTP_PORT

RUN ls -la

# Run the app binary
CMD ["./app"]