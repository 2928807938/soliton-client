# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.work and all module files
COPY go.work ./
COPY bom/go.mod ./bom/
COPY share/go.mod ./share/
COPY user/go.mod ./user/
COPY user/domain/go.mod ./user/domain/
COPY user/infrastructure/go.mod ./user/infrastructure/
COPY api/go.mod ./api/
COPY api/user-api/go.mod ./api/user-api/
COPY cmd/api/go.mod ./cmd/api/

# Download dependencies
RUN go work sync

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
