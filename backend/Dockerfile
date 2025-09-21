FROM golang:1.22-alpine AS build
RUN apk add --no-cache build-base
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server /app/server
COPY --from=build /app/migrations /app/migrations
COPY .env .env
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["/app/server"]
