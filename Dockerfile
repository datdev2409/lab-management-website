FROM golang:1.23.9-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GO_ENV=production go build -o /app/bin/main /app/cmd/api/main.go

# Deploy the application binary into a lean image
FROM alpine:latest

WORKDIR /

COPY --from=builder /app/bin/main /main
COPY --from=builder /app/.env.production /.env.production

EXPOSE 8080

ENV GO_ENV=production

ENTRYPOINT ["/main"]

