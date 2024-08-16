FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache ca-certificates
COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
