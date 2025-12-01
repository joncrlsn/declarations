# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
# Enable Go modules and download dependencies
COPY go.mod ./
RUN go mod tidy || true
# Copy the rest of the source
COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Runtime stage
FROM alpine:3.20
WORKDIR /app
# Add CA certificates (for any outbound HTTPS calls)
RUN apk --no-cache add ca-certificates
# Copy binary from builder
COPY --from=builder /app/server /app/server
COPY declarations.txt /app/declarations.txt
#COPY .api-token /app/.api-token

EXPOSE 8080
# Run the server
CMD ["/app/server"]

