FROM golang:latest AS builder

WORKDIR /app

# important step: download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the app in the container
# WARNING: be sure to use .gitignore and .dockerignore files
COPY . . 

# might need env vars

# Build with verbose output
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

# second stage for smaller image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

# Make it executable
RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]


