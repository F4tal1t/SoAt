FROM golang:1.24 AS builder

WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o SoAt cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/SoAt /app/SoAt

EXPOSE 3000

CMD ["./SoAt"]