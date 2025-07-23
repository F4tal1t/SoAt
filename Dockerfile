FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o SoAt cmd/main.go

# Build the admin creation utility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o create_admin cmd/create_admin/main.go

# Build the role column utility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o add_role_column cmd/add_role_column/main.go

# Use a smaller base image for the final stage
FROM alpine:3.20

# Add necessary dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=UTC

# Create a non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/SoAt /app/SoAt
COPY --from=builder /app/create_admin /app/create_admin
COPY --from=builder /app/add_role_column /app/add_role_column

# Set ownership
RUN chown -R appuser:appuser /app

# Use the non-root user
USER appuser

EXPOSE 3000

CMD ["./SoAt"]