FROM golang:1.22-alpine
RUN apk add --no-cache nodejs npm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN npm install
RUN npm run build:css

RUN go build -o main cmd/gitea-committer/main.go


EXPOSE 8080

CMD ["./main"]
