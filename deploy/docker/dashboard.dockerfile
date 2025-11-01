FROM node:22-alpine AS client-builder
WORKDIR /app/client
RUN apk add --no-cache python3 make g++
COPY dashboard/client/package*.json ./
RUN npm install
COPY dashboard/client/ .
RUN npm run build  # Outputs to ../server/dist

FROM golang:1.25-alpine AS server-builder
WORKDIR /app/server
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY dashboard/server/ .
COPY lib ./lib
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o dashboard .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=client-builder /app/server/dist ./dist
COPY --from=server-builder /app/server/dashboard .
EXPOSE 8080
CMD ["./dashboard"]