FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o main cmd/gitea-committer/main.go


EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
