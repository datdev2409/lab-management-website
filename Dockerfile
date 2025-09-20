FROM golang:1.24.7-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate
RUN mkdir -p /app/reports

RUN go build -o /app/bin/main /app/cmd/api/main.go

FROM node:alpine3.22 as esbuild

WORKDIR /app

COPY internal/templates/scripts ./internal/templates/scripts

RUN npm install -g esbuild

RUN esbuild ./internal/templates/scripts/ --bundle --minify --outfile=./static/index.js

# Deploy the application binary into a lean image
FROM scratch

WORKDIR /

# Copy TLS certificates to allow TLS traffic
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/bin/main /main
COPY --from=builder /app/templates /templates
COPY --from=builder /app/reports /reports
COPY --from=esbuild /app/static/index.js /static/index.js

# Create empty reports directory
EXPOSE 9000

ENTRYPOINT ["/main"]

