FROM node:22-alpine AS client-builder
WORKDIR /app/client
RUN apk add --no-cache python3 make g++
COPY client/package*.json ./
RUN npm install
COPY client/ .
RUN npm run build  # Outputs to ../cmd/dashboard/dist

FROM golang:1.25-alpine AS server-builder
WORKDIR /app/server
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
COPY --from=client-builder /app/cmd/dashboard/dist ./cmd/dashboard/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o dashboard cmd/dashboard/*.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=server-builder /app/server/dashboard .
EXPOSE 8080
CMD ["./dashboard"]