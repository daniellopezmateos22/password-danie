FROM golang:1.22-alpine AS build
RUN apk add --no-cache build-base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server /app/server
COPY --from=build /app/migrations /app/migrations
RUN mkdir -p /data
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
