# ---- build ----
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd

# ---- runtime ----
FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server .
COPY .env .env
EXPOSE 8080
CMD ["./server"]
