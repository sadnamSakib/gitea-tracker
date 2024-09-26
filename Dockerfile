FROM ghcr.io/a-h/templ:latest AS generate-stage

WORKDIR /app


COPY --chown=65532:65532 . /app


RUN ["templ", "generate"]

FROM golang:1.22-alpine AS build-stage


RUN apk add --no-cache nodejs npm build-base git tzdata


WORKDIR /app

COPY go.mod go.sum ./


RUN go mod download


COPY --from=generate-stage /app /app


RUN npm install
RUN npm run build:css


RUN go build -o /app/main ./cmd/gitea-committer/main.go


FROM alpine:3.18


RUN apk add --no-cache ca-certificates tzdata


WORKDIR /app


COPY --from=build-stage /app/main /app/main
COPY --from=build-stage /app/web /app/web


ENV TZ=Asia/Dhaka


EXPOSE 8080


CMD ["./main"]
