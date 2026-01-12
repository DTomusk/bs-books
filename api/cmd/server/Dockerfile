# Stage 1: Build binaries 
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /bin/bs-books-api ./cmd/server


# -------- test image --------
FROM builder AS test
CMD ["go", "test", "./internal/...", "-v"]

# -------- production image --------

FROM alpine:3.22.2

ENV ENV=docker

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/bs-books-api /usr/local/bin/

CMD ["/usr/local/bin/bs-books-api"]