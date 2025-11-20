# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -l .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

# Final stage – distroless is very popular on Cloud Run
#FROM gcr.io/distroless/static-debian12
FROM alpine:3.22
COPY --from=builder server /app/declarations.txt /
RUN groupadd --gid 1000 nonrootgroup && useradd -u 1000 -g nonrootgroup nonroot
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
