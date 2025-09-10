#!/bin/bash

# Pre-commit script to run tests and checks locally
# This mimics what the CI pipeline will do

set -e

echo "🚀 Running pre-commit checks..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi

# Check if templ is installed
if ! command -v templ &> /dev/null; then
    echo "📦 Installing templ..."
    go install github.com/a-h/templ/cmd/templ@latest
fi

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "📦 Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo "✅ All dependencies installed"

# Generate templ files
echo "🔧 Generating templ files..."
templ generate

# Format code
echo "🎨 Formatting code..."
go fmt ./...

# Download and verify dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod verify

# Run linter
echo "🔍 Running linter..."
golangci-lint run

# Run go vet
echo "🔍 Running go vet..."
go vet ./...

# Check if gofmt would make changes
echo "🔍 Checking code formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "❌ The following files are not formatted:"
    gofmt -s -l .
    echo "Run 'go fmt ./...' to fix formatting issues"
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
go test -v -race ./...

# Build the application
echo "🔨 Building application..."
mkdir -p bin
GOOS=linux GOARCH=amd64 go build -o bin/main cmd/api/main.go

echo "✅ All checks passed! You're ready to push."
