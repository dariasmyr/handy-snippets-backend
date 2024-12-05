FROM golang:alpine AS build
RUN apk update && apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/server cmd/server.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/server /app/server
CMD ["./server"]
