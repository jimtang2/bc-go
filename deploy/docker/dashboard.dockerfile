# webhook$$bc-go;deploy/docker/dashboard.dockerfile;grok$$
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dashboard dashboard/server/*.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/dashboard .
COPY dashboard/server/dist ./dist
CMD ["./dashboard"]