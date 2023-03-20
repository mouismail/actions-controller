# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19 AS builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.mod and go.sum files to the workspace.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code to the workspace.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go test ./...

# Use a minimal alpine image to run the binary.
# https://hub.docker.com/_/alpine
FROM alpine:latest

# Copy the binary from the builder stage.
COPY --from=builder /app .

# Run the binary.
CMD ["./app"]
