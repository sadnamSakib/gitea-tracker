# Stage 1: Build the Go application
FROM golang:1.22-alpine AS build-stage

# Install dependencies
RUN apk add --no-cache nodejs npm build-base git

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire application code
COPY . .

# Install npm dependencies and build CSS
RUN npm install
RUN npm run build:css

# Build the Go application
RUN  go build -o /app/main ./cmd/gitea-committer/main.go

# Stage 2: Generate templ files
FROM ghcr.io/a-h/templ:latest AS generate-stage

# Set working directory and copy the app source code
WORKDIR /app
COPY --chown=65532:65532 . .

# Generate templ output
RUN templ generate

# Stage 3: Create a lightweight final image
FROM alpine:3.18

# Install required packages
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy the Go binary and web assets from the build-stage
COPY --from=build-stage /app/main /app/main
COPY --from=build-stage /app/web /app/web

# Expose the application port
EXPOSE 8080

# Run the Go application
CMD ["./main"]
