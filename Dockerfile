FROM golang:1.23.10-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GO_ENV=production go build -o /app/bin/main /app/cmd/api/main.go

# Deploy the application binary into a lean image
FROM scratch

WORKDIR /

# Copy TLS certificates to allow TLS traffic
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/bin/main /main
COPY --from=builder /app/.env.production /.env.production
COPY --from=builder /app/templates /templates

EXPOSE 8080

ENV GO_ENV=production

ENTRYPOINT ["/main"]

