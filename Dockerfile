FROM golang:1.23.10-alpine as builder

ARG ENV
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate
RUN mkdir -p /app/reports

RUN GO_ENV=${ENV} go build -o /app/bin/main /app/cmd/api/main.go

# Deploy the application binary into a lean image
FROM scratch

ARG ENV
WORKDIR /

# Copy TLS certificates to allow TLS traffic
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/bin/main /main
COPY --from=builder /app/templates /templates
COPY --from=builder /app/reports /reports

# Create empty reports directory
EXPOSE 8080

ENV GO_ENV=${ENV}

ENTRYPOINT ["/main"]

