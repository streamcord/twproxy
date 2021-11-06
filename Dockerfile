FROM golang:1.16-alpine

WORKDIR /app/
COPY ./ /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o twproxy
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=0 /app/twproxy /app/

ENV GIN_MODE=release

ENTRYPOINT ["/app/twproxy"]