FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the app binary with ldflags
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG COMMIT_ID=local
ARG HTTP_PORT=3000
ARG GHES_APP_WEBHOOK_SECRET=development
ARG GHES_APP_PRIVATE_KEY=/app/test/actions-control.2023-03-30.private-key.pem

ENV GHES_APP_WEBHOOK_SECRET=$GHES_APP_WEBHOOK_SECRET
ENV GHES_APP_PRIVATE_KEY=$GHES_APP_PRIVATE_KEY

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-w -s \
    -X main.version=$VERSION \
    -X main.buildTime=$BUILD_TIME \
    -X main.commitID=$COMMIT_ID"

# Use alpine:latest as base image for final stage
FROM alpine:latest

# Copy binary from builder stage to current stage
COPY --from=builder /app /app

# Expose port 3000 to the outside world
EXPOSE $HTTP_PORT

RUN chmod +x ./app
RUN ls /app

# Run the app binary
CMD [ "/app/actions-rollout-app" ]