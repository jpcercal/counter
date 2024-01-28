FROM golang:1.21 AS builder

WORKDIR /src

# Install dependencies
COPY go.mod go.sum /
RUN go mod download

# Copy source files
COPY . .
RUN CGO_ENABLED=0 go build -o main .

# production-ready image
FROM alpine AS production
RUN apk update && apk upgrade && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy binary
COPY --from=builder /src/main .

# Expose port
EXPOSE 3000

# Set the binary as the entrypoint of the container
CMD ["./main", "serve"]
