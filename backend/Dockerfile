FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:3.20
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
