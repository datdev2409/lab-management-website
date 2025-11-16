# Lab Management System

A comprehensive web-based laboratory information management system for **Anh Quan Laboratory**. This system streamlines laboratory operations by managing patient records, test definitions, test combos, and providing tracking and comparison features for test results.

[![Go Version](https://img.shields.io/badge/Go-1.24.7-blue.svg)](https://golang.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-v2-green.svg)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📋 Table of Contents

- [Features](#✨-features)
  - [Core Functionality](#core-functionality)
  - [Advanced Features](#advanced-features)
- [Technology Stack](#🛠-technology-stack)
  - [Backend](#backend)
  - [Frontend](#frontend)
  - [DevOps & Tools](#devops--tools)
- [Architecture](#🏗-architecture)
- [Local development (quickstart)](#🚀-local-development-quickstart)
- [Deployment](#🚢-deployment-github-actions-cd---recommended)
- [Project Structure](#📁-project-structure)
- [Contributing](#🤝-contributing)
- [License](#📄-license)
- [Contact & Support](#📞-contact--support)

## ✨ Features

### Core Functionality

- **Authentication**: JWT-based with HTTP-only cookies authentication.
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
- **Bootstrap**: ResponsiveUI framework
- **ESBuild**: Fast JavaScript bundler

### DevOps & Tools

- **Docker & Docker Compose**: Containerization
- **Air**: Live reload for Go development
- **Gotenberg**: PDF generation service
- **Traefik**: Reverse proxy (production)

## 🏗 Architecture

The system follows **Domain-Driven Design** principles with clear separation of concerns:

```
      ┌─────────────────┐
      │  Web Browser    │
      │  (HTMX/Alpine)  │
      └────┬────────────┘
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

## 🚀 Local development (quickstart)

A minimal, fast path to get the app running locally.

1. Prerequisites

   - Go 1.24.7+, Node.js (optional for tests/ESBuild), Docker (optional)

2. Clone & dependencies

   ```bash
   git clone https://github.com/datdev2409/lab-management-website.git
   cd lab-management-website
   go mod download
   npm install         # optional: needed for Playwright / ESBuild
   ```

3. Create configuration

   ```bash
   cp .env.example .env
   # Edit .env and update values: MONGODB_URI, JWT_SECRET, GOTENBERG_URL, etc.
   ```

4. Start local services (recommended)

   ```bash
   docker-compose up -d
   # starts MongoDB and Gotenberg; skip if you have services already
   ```

5. Install small dev tools (one-time)

   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   go install github.com/air-verse/air@latest   # optional: hot reload
   npm install -g esbuild                       # optional: JS bundling
   ```

6. Generate templates and build assets

   ```bash
   templ generate
   # if you changed frontend assets:
   npm run build   # or run make live/esbuild during development
   ```

7. Start the app (recommended: 3-way hot-reload)

   ```bash
   make live          # runs server, templ watcher, and esbuild concurrently
   # OR start only the server:
   make live/server
   # OR run directly:
   go run ./cmd/api
   ```

8. Access the app

   - Templ proxy (auto reload when templ file changes; recommended during development): http://localhost:7331
   - Direct server: http://localhost:9000

9. Testing & quick commands
   - Run E2E (UI): make e2e
   - Run E2E (CI/headless): make e2e-ci
   - Build dev binary: go build -o bin/main cmd/api/main.go
   - Production build (Linux): make build
   - Docker image: docker build -t lab-management:latest .

Tips

- Use docker-compose to simplify local dependencies; stop services with docker-compose down.
- Use make live for the fastest feedback loop (templ + server + esbuild).
- Keep .env out of git; use .env.production for staging/prod configs.

## 🚢 Deployment (GitHub Actions CD - recommended)

This project uses the GitHub Actions CD pipeline (.github/workflows/cd.yml) for manual deployments. The workflow is triggered with workflow_dispatch and accepts two inputs:

- image_tag (optional): commit SHA or tag. If omitted, the workflow uses the branch's latest commit SHA.
- environment (required): either `stg` or `prod`.

Recommended deployment steps

1. Verify the CI pipeline is green with the target commit / branch
2. Ensure repository variables & secrets are set:

   - Repository variables: SSH_HOST, SSH_USER, DOCKER_USERNAME
   - Repository secret: SSH_PRIVATE_KEY
   - These are referenced by the workflow and used by the SSH action.

3. Trigger the CD workflow (via GitHub webiste or gh CLI):

   ```bash
   # use commit SHA or a tag for image_tag; omit image_tag to deploy workflow's github.sha
   gh workflow run cd.yml -f image_tag=962bb4e471b827a2eb9f9706912e3aa69dbb1a36 -f environment=stg
   ```

4. Optional: Create release/tag after a successful deployment
   ```bash
   git tag -a v1.0.1 <commit-sha> -m "Release v1.0.1"
   git push origin v1.0.1
   gh release create v1.0.1 --title "Release v1.0.1" --notes "Deployed to stg"
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
├── tests/                          # Playwright E2E tests
├── deploy/                         # Deployment configurations
├── docs/                           # Documentation & Bruno Collection
├── scripts/                        # Utility scripts
├── Makefile                        # Build and development commands
├── docker-compose.yaml             # Local development services
├── Dockerfile                      # Multi-stage production build
└── README.md                       # This file
```

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
- Update e2e tests if there is any changes in UI
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

---

**Built with ❤️ for Anh Quan Laboratory**
