# CI/CD Pipeline Documentation

This document describes the CI/CD pipeline setup for the Lab Admin Go application using GitHub Actions.

## Overview

The CI/CD pipeline consists of multiple workflows that handle different aspects of the development lifecycle:

1. **Continuous Integration (CI)** - `ci.yml`
2. **Continuous Deployment (CD)** - `cd.yml`
3. **Security Scanning** - `security.yml`
4. **Dependency Updates** - `dependency-updates.yml`
5. **Release Management** - `release.yml`

## Workflows

### 1. CI Pipeline (`ci.yml`)

**Triggers:**

- Push to `main` and `develop` branches
- Pull requests to `main` and `develop` branches

**Jobs:**

- **Test**: Runs unit tests with MongoDB service
- **Lint**: Code quality checks using golangci-lint
- **Security**: Security scanning with Gosec
- **Build**: Builds the application binary
- **E2E Tests**: End-to-end testing with Cypress

**Features:**

- Go module caching for faster builds
- Templ file generation
- Code coverage reporting to Codecov
- Artifact uploads for debugging

### 2. CD Pipeline (`cd.yml`)

**Triggers:**

- Push to `main` branch (production deployment)
- Push to `develop` branch (staging deployment)
- Manual workflow dispatch
- Git tags starting with `v*`

**Jobs:**

- **Build and Push**: Creates Docker images and pushes to GitHub Container Registry
- **Deploy Staging**: Deploys to staging environment (develop branch)
- **Deploy Production**: Deploys to production environment (main branch/tags)
- **Notify**: Sends deployment notifications

**Features:**

- Multi-platform Docker builds (AMD64 + ARM64)
- Rolling deployments with health checks
- Automated rollback on failure
- Environment-specific configurations

### 3. Security Pipeline (`security.yml`)

**Triggers:**

- Daily schedule (2 AM UTC)
- Push to `main` and `develop` branches
- Pull requests to `main`
- Manual workflow dispatch

**Jobs:**

- **CodeQL Analysis**: Static analysis for Go and JavaScript
- **Docker Security**: Container vulnerability scanning with Trivy
- **Dependency Check**: Go module vulnerability scanning
- **Secrets Scan**: Secret detection with TruffleHog

### 4. Dependency Updates (`dependency-updates.yml`)

**Triggers:**

- Weekly schedule (Mondays at 9 AM UTC)
- Manual workflow dispatch

**Jobs:**

- **Update Go Dependencies**: Updates Go modules and creates PR
- **Update Docker Images**: Updates base Docker images and creates PR

### 5. Release Management (`release.yml`)

**Triggers:**

- Manual workflow dispatch with version input

**Jobs:**

- **Validate**: Validates version format and runs tests
- **Create Release**: Creates Git tag and GitHub release
- **Deploy Production**: Triggers production deployment

## Setup Instructions

### 1. Repository Secrets

Configure the following secrets in your GitHub repository settings:

#### Staging Environment

- `STAGING_HOST`: Staging server hostname/IP
- `STAGING_USER`: SSH username for staging server
- `STAGING_SSH_KEY`: SSH private key for staging server
- `STAGING_PORT`: SSH port (optional, defaults to 22)

#### Production Environment

- `PRODUCTION_HOST`: Production server hostname/IP
- `PRODUCTION_USER`: SSH username for production server
- `PRODUCTION_SSH_KEY`: SSH private key for production server
- `PRODUCTION_PORT`: SSH port (optional, defaults to 22)

### 2. Environment Setup

#### Staging Server Setup

```bash
# Create application directory
sudo mkdir -p /opt/lab-admin-go
cd /opt/lab-admin-go

# Copy docker-compose.staging.yaml to server
# Ensure Docker and Docker Compose are installed
```

#### Production Server Setup

```bash
# Create application directory
sudo mkdir -p /opt/lab-admin-go
cd /opt/lab-admin-go

# Copy docker-compose.prod.yaml to server
# Ensure Docker and Docker Compose are installed
```

### 3. GitHub Container Registry

The pipeline automatically builds and pushes Docker images to GitHub Container Registry (ghcr.io). No additional setup is required as it uses the `GITHUB_TOKEN`.

### 4. Branch Protection

Configure branch protection rules for `main` and `develop` branches:

1. Go to Settings > Branches
2. Add protection rules for `main` and `develop`
3. Enable:
   - Require status checks to pass before merging
   - Require branches to be up to date before merging
   - Require pull request reviews before merging
   - Dismiss stale PR approvals when new commits are pushed

## Usage

### Development Workflow

1. **Feature Development**:

   ```bash
   git checkout develop
   git checkout -b feature/your-feature
   # Make changes
   git commit -m "feat: add new feature"
   git push origin feature/your-feature
   ```

2. **Create Pull Request**:

   - Create PR from feature branch to `develop`
   - CI pipeline runs automatically
   - Wait for all checks to pass
   - Request code review
   - Merge after approval

3. **Staging Deployment**:

   - When PR is merged to `develop`, automatic deployment to staging occurs
   - Test your changes in staging environment

4. **Production Release**:
   - Create PR from `develop` to `main`
   - After merge, create a release using the Release workflow
   - Production deployment happens automatically

### Manual Deployment

You can trigger manual deployments using the workflow dispatch feature:

1. Go to Actions tab in GitHub
2. Select "CD Pipeline" workflow
3. Click "Run workflow"
4. Choose environment (staging/production)
5. Click "Run workflow"

### Creating Releases

1. Go to Actions tab in GitHub
2. Select "Release" workflow
3. Click "Run workflow"
4. Enter version (e.g., v1.0.0)
5. Add release notes (optional)
6. Click "Run workflow"

## Monitoring and Debugging

### Build Logs

- All workflow runs are logged in the Actions tab
- Click on any workflow run to see detailed logs
- Download artifacts for debugging (build outputs, test results, etc.)

### Health Checks

- Production deployment includes health checks
- Application provides `/health` endpoint
- Deployment fails if health checks don't pass

### Rollback Procedure

If a deployment fails or issues are discovered:

1. **Automatic Rollback**: Health checks will prevent bad deployments
2. **Manual Rollback**:
   ```bash
   # On the server
   docker-compose -f docker-compose.prod.yaml down
   # Change image tag to previous version in compose file
   docker-compose -f docker-compose.prod.yaml up -d
   ```

## Troubleshooting

### Common Issues

1. **Build Failures**:

   - Check Go module compatibility
   - Ensure templ files generate correctly
   - Verify test dependencies

2. **Deployment Failures**:

   - Check SSH connectivity to servers
   - Verify Docker service is running on target servers
   - Check disk space on target servers

3. **Health Check Failures**:
   - Verify application starts correctly
   - Check environment variables
   - Review application logs

### Debug Commands

```bash
# Check application logs
docker-compose -f docker-compose.prod.yaml logs service

# Check container status
docker-compose -f docker-compose.prod.yaml ps

# Restart services
docker-compose -f docker-compose.prod.yaml restart service

# Check health endpoint
curl -f http://localhost:8888/health
```

## Best Practices

1. **Always create feature branches** from `develop`
2. **Write tests** for new features
3. **Keep commits atomic** and well-described
4. **Review code** before merging
5. **Test in staging** before production release
6. **Monitor deployments** and application health
7. **Keep dependencies updated** (automated weekly)
8. **Follow semantic versioning** for releases

## Configuration Files

- `.github/workflows/`: GitHub Actions workflow files
- `.golangci.yml`: Linting configuration
- `docker-compose.staging.yaml`: Staging environment configuration
- `docker-compose.prod.yaml`: Production environment configuration
- `Dockerfile`: Application container configuration

## Contact

For questions about the CI/CD pipeline, please create an issue or contact the development team.
