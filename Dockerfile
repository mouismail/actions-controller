FROM golang:1.19.8-alpine3.17 as builder

WORKDIR /app

VOLUME /app/keys

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/actions-rollout-app /app
COPY --from=builder /app/actions-controller.yaml /app
COPY --from=builder /app/. /app

ARG GHES_APP_PRIVATE_KEY
ARG GHES_APP_WEBHOOK_SECRET

ENV GHES_APP_PRIVATE_KEY=$GHES_APP_PRIVATE_KEY
ENV GHES_APP_WEBHOOK_SECRET=$GHES_APP_WEBHOOK_SECRET

EXPOSE 3000
CMD ["/app/actions-rollout-app", "-c", "/app/actions-controller.yaml"]