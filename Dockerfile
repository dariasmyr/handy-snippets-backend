FROM golang:alpine AS build
RUN apk update && apk add --no-cache make gcc musl-dev
ENV CGO_ENABLED=1
WORKDIR /app
COPY . .
RUN make build
FROM alpine:latest
RUN apk update && apk add --no-cache make gcc musl-dev
WORKDIR /app
COPY --from=build /app/build/server /app/server
CMD ["./server"]
