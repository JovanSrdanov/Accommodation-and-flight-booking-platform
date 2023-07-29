# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./reservation_service/go.mod ./reservation_service/go.sum ./

# Copy the local dependency
COPY /common ../common
COPY /authorization_service ../authorization_service
# Copy the cert folder from the build context to the container
COPY ./authorization_service/cert ./cert

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./reservation_service/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the cert folder from the first stage to the second stage
COPY --from=builder /app/cert ./cert

# Expose port 8000 to the outside world
EXPOSE 8000
RUN chmod +x ./main
# Command to run the executable
CMD ["./main"]