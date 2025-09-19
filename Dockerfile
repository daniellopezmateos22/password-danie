# ---- build ----
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd

# ---- runtime ----
FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
