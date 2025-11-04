# Lab Management System

A comprehensive web-based laboratory information management system for **Anh Quan Laboratory**. This system streamlines laboratory operations by managing patient records, test definitions, test combos, and providing tracking and comparison features for test results.

[![Go Version](https://img.shields.io/badge/Go-1.24.7-blue.svg)](https://golang.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-v2-green.svg)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📋 Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development](#development)
- [Testing](#testing)
- [Building](#building)
- [Deployment](#deployment)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

### Core Functionality
- **Patient Management**: Create, update, and search patient records with auto-complete functionality
- **Test Management**: Define test types with normal ranges, pricing, and units
- **Combo Management**: Create reusable test packages for common test combinations
- **Record Management**: Complete workflow for creating and managing lab test records
- **Status Tracking**: Track record status from pending to completed
- **Result Comparison**: Compare test results across multiple records to monitor patient progress
- **Report Generation**: Generate multiple report types (billing, results, signed results, PDF reports, tracking reports)

### Advanced Features
- **Abnormal Test Detection**: Automatic detection of abnormal test results based on defined ranges
- **Tracking & Comparison**: Create tracking configurations to compare test results over time
- **Unsaved Changes Warning**: Alerts users before leaving pages with unsaved modifications
- **Revenue Reporting**: Generate comprehensive revenue reports for laboratory operations

## 🛠 Technology Stack

### Backend
- **Go 1.24.7**: High-performance backend API
- **Chi Router**: Lightweight HTTP router
- **MongoDB v2**: NoSQL database for flexible data storage
- **JWT**: Secure authentication with HTTP-only cookies
- **Zap**: Structured logging

### Frontend
- **HTMX**: Modern dynamic web interactions (being migrated to Alpine.js)
- **Alpine.js**: Lightweight JavaScript framework for interactivity
- **Templ**: Type-safe Go templating engine
- **Bootstrap**: Responsive UI framework
- **ESBuild**: Fast JavaScript bundler

### DevOps & Tools
- **Docker & Docker Compose**: Containerization
- **Air**: Live reload for Go development
- **Gotenberg**: PDF generation service
- **Excelize**: Excel file generation
- **Traefik**: Reverse proxy (production)
- **systemd**: Service management (production)

## 🏗 Architecture

The system follows **Domain-Driven Design** principles with clear separation of concerns:

```
┌─────────────────┐
│  Web Browser    │
│  (HTMX/Alpine)  │
└────────┬────────┘
         │
    ┌────▼────────────────┐
    │   HTTP Handlers     │
    │  (Controllers)      │
    └────┬────────────────┘
         │
    ┌────▼────────────────┐
    │   Business Logic    │
    │   (Models/DTOs)     │
    └────┬────────────────┘
         │
    ┌────▼────────────────┐
    │   Storage Layer     │
    │   (Data Access)     │
    └────┬────────────────┘
         │
    ┌────▼────────────────┐
    │   MongoDB           │
    └─────────────────────┘
```

### Key Components
- **Authentication**: JWT-based with HTTP-only cookies
- **Templating**: Server-side rendering with Templ
- **Report Generation**: Excel files with Excelize, PDF conversion via Gotenberg
- **Logging**: Structured logging with Zap
- **Database**: MongoDB with custom ID generation (`type_randomstring` pattern)

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.24.7 or higher ([Download](https://golang.org/dl/))
- **MongoDB**: Version 4.0 or higher ([Download](https://www.mongodb.com/try/download/community))
- **Node.js**: For Cypress testing and ESBuild ([Download](https://nodejs.org/))
- **Docker & Docker Compose**: For containerized deployment (optional) ([Download](https://www.docker.com/))
- **Templ**: Go templating engine
  ```bash
  go install github.com/a-h/templ/cmd/templ@latest
  ```
- **Air**: Live reload tool for development
  ```bash
  go install github.com/air-verse/air@latest
  ```
- **golangci-lint**: For code linting
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```
- **ESBuild**: JavaScript bundler
  ```bash
  npm install -g esbuild
  ```

## 🚀 Installation

### 1. Clone the Repository
```bash
git clone https://github.com/datdev2409/lab-management-website.git
cd lab-management-website
```

### 2. Install Go Dependencies
```bash
go mod download
go mod verify
```

### 3. Install Testing Dependencies (Optional)
```bash
cd tests
npm install
cd ..
```

### 4. Start MongoDB
Using Docker Compose (recommended):
```bash
docker-compose up -d
```

This starts:
- MongoDB on port 27017 (credentials: root/password123)
- Gotenberg PDF service on port 3000

Or start MongoDB manually:
```bash
mongod --dbpath /your/data/path
```

### 5. Generate Templ Files
```bash
templ generate
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Server Configuration
SERVER_PORT=9000
ENV=local

# MongoDB Configuration
MONGODB_URI=mongodb://root:password123@localhost:27017/labadmin?authSource=admin

# JWT Configuration
JWT_SECRET=your-secret-key-here

# Gotenberg Service (PDF Generation)
GOTENBERG_URL=http://localhost:3000
```

### Environment Files by Stage
- **Local Development**: `.env` or `.env.local`
- **Production**: `.env.production`

## 💻 Development

### Live Development Mode

The project supports **3-way hot reload** for rapid development:

```bash
make live
```

This command runs three concurrent processes:
1. **Air**: Go hot reload - restarts server on `.go` file changes
2. **Templ**: Template generation with live reload on `.templ` file changes
3. **ESBuild**: JavaScript bundling and minification on script changes

Access the application at: `http://localhost:7331` (Templ proxy) or `http://localhost:9000` (direct)

### Individual Development Commands

Run specific development processes separately:

```bash
# Go server with hot reload
make live/server

# Templ generation with watch mode
make live/templ

# JavaScript bundling with watch mode
make live/esbuild

# Start only MongoDB and Gotenberg
make env
```

### Code Quality

Run pre-commit checks before pushing:

```bash
./scripts/pre-commit.sh
```

This script runs:
- Templ file generation
- Code formatting (`go fmt`)
- Dependency verification
- Linting (`golangci-lint`)
- Static analysis (`go vet`)
- Format checking
- Unit tests
- Build verification

## 🧪 Testing

### Unit Tests

Run Go unit tests:

```bash
go test -v ./...
```

Run tests with race detection:

```bash
go test -v -race ./...
```

### End-to-End Tests

Run Cypress E2E tests:

```bash
cd tests
npm test        # Run headless
npm run cy:open # Open Cypress UI
```

## 🔨 Building

### Development Build

```bash
go build -o bin/main cmd/api/main.go
```

### Production Build

Build for Linux (x86_64):

```bash
make build
```

This creates an optimized binary at `bin/main`.

### Docker Build

Build Docker image:

```bash
docker build -t lab-management:latest .
```

The Dockerfile uses a multi-stage build:
1. **Builder stage**: Compiles Go application with Templ generation
2. **ESBuild stage**: Bundles JavaScript files
3. **Runtime stage**: Minimal Alpine image with compiled binary

## 🚢 Deployment

### Local Docker Deployment

```bash
docker-compose -f docker-compose.local.yaml up
```

### Production Deployment

The system uses **systemd** for service management in production. See [deploy/deploy.md](deploy/deploy.md) for detailed instructions.

#### Quick Deployment Steps:

1. **Build the application**:
   ```bash
   make build
   ```

2. **Copy binary to server**:
   ```bash
   scp bin/main user@server:/home/user/
   ```

3. **Set up environment**:
   ```bash
   # On server
   cp .env.example .env.production
   # Edit .env.production with production values
   ```

4. **Install and start systemd service**:
   ```bash
   sudo cp deploy/goweb.service /lib/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable goweb
   sudo systemctl start goweb
   ```

5. **Check service status**:
   ```bash
   sudo systemctl status goweb
   journalctl -u goweb -f
   ```

### Docker Compose Deployment

For containerized deployment with base infrastructure:

```bash
# Start base services (MongoDB, Traefik, etc.)
make start-base ENV=production

# Start application
make start-app ENV=production DOCKER_USERNAME=your-username DOCKER_TAG=latest

# Stop services
make stop-app ENV=production
make stop-base ENV=production
```

## 📁 Project Structure

```
lab-management-website/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── auth/                       # JWT authentication
│   ├── db/                         # Database connection
│   ├── handlers/                   # HTTP handlers (controllers)
│   │   ├── handler.go             # Router setup
│   │   ├── auth_handler.go        # Authentication endpoints
│   │   ├── patient_handler.go     # Patient management
│   │   ├── record_handler.go      # Lab record management
│   │   ├── test_handler.go        # Test definition management
│   │   ├── combo_handler.go       # Test combo management
│   │   ├── tracking_handler.go    # Tracking/comparison
│   │   └── middlewares.go         # JWT auth, logging
│   ├── models/                     # Data models and DTOs
│   ├── storage/                    # Data access layer
│   │   ├── storage.go             # Storage interface
│   │   ├── base.go                # Generic MongoDB operations
│   │   └── *_storage.go           # Entity-specific CRUD
│   ├── sheets/                     # Excel/PDF generation
│   ├── templates/                  # Templ templates
│   │   ├── pages/                 # Full page templates
│   │   ├── partials/              # Reusable components
│   │   └── scripts/               # JavaScript files
│   └── logger/                     # Structured logging
├── templates/                      # Excel templates
├── reports/                        # Generated reports (gitignored)
├── tests/                          # Cypress E2E tests
├── deploy/                         # Deployment configurations
├── docs/                           # Additional documentation
├── scripts/                        # Utility scripts
├── Makefile                        # Build and development commands
├── docker-compose.yaml             # Local development services
├── Dockerfile                      # Multi-stage production build
└── README.md                       # This file
```

### Domain Models

- **Patient**: Patient information (`patient_*` ID)
- **Test**: Test definitions with ranges and pricing (`test_*` ID)
- **Combo**: Test packages (`combo_*` ID)
- **Record**: Lab test records with results (`record_*` ID)
- **Tracking**: Comparison configurations (`tracking_*` ID)
- **User**: Authentication and authorization (`user_*` ID)

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** following the coding conventions
4. **Run pre-commit checks**: `./scripts/pre-commit.sh`
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Coding Conventions

- Follow Go standard formatting (`gofmt`, `go vet`)
- Use meaningful variable and function names
- Write unit tests for new functionality
- Document exported functions and types
- Use structured logging with Zap
- Follow the existing project structure

### Pre-commit Checklist

Before committing, ensure:
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No linting errors (`golangci-lint run`)
- [ ] All tests pass (`go test ./...`)
- [ ] Templ files generated (`templ generate`)
- [ ] Application builds successfully (`make build`)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact & Support

- **Project Owner**: [datdev2409](https://github.com/datdev2409)
- **Repository**: [lab-management-website](https://github.com/datdev2409/lab-management-website)
- **Issues**: [GitHub Issues](https://github.com/datdev2409/lab-management-website/issues)

For detailed implementation documentation, see the [docs](docs/) directory:
- [CD Pipeline Documentation](docs/CD_PIPELINE.md)
- [Unsaved Changes Feature](docs/README-unsaved-changes.md)
- [Deployment Guide](deploy/deploy.md)

---

**Built with ❤️ for Anh Quan Laboratory**
