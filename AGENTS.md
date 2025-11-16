# AI Agent Instructions for Lab Management System

This document provides guidelines for AI agents working on this project, including development workflows, testing, production deployment and GitHub workflows.

## Dev Environments

## Production Deployment Workflows

### Overview

This section documents the complete deployment workflow for the Lab Management System, from code verification through release tagging. All deployments must be manually triggered—no automatic deployments.

### Pre-Deployment Verification

- **Gate deployments**: Always verify all GitHub Actions status checks for the target commit are green before deploying. Use GitHub MCP to query the commit check suite and block deployment if any required check is failing or pending.
- **Manual-only**: Do not run production/CD workflows automatically. Trigger production deployments only when explicitly requested by the user.
- **Confirmation**: Ask user to confirm target commit SHA/branch and environment (stg | prod) before proceeding.

### Branch Deployments

- When asked to deploy a branch, fetch the branch tip commit SHA and ensure its checks are passing.
- Deploy the latest commit on that branch only.
- Example: `deploy latest commit on develop to stg` → fetch develop tip → verify checks → deploy SHA.

### Triggering CD Workflow

- Use `gh workflow run cd.yml` via terminal to dispatch the CD workflow with inputs:
  - `image_tag`: commit SHA or image tag to deploy
  - `environment`: stg | prod (as requested)
- Example command: `gh workflow run cd.yml -f image_tag=962bb4e471b827a2eb9f9706912e3aa69dbb1a36 -f environment=stg`
- Workflow definition: `.github/workflows/cd.yml` (uses `workflow_dispatch` event)

### Workflow Monitoring

- After triggering, poll the workflow run status via GitHub API: `curl -s -H "Accept: application/vnd.github.v3+json" "https://api.github.com/repos/{owner}/{repo}/actions/workflows/cd.yml/runs?per_page=1" | jq '.workflow_runs[0] | {id, status, conclusion}'`
- Expected completion time: ~2-3 minutes
- Success indicators: `status: "completed"` and `conclusion: "success"`
- If workflow fails: Report failure with run URL and logs to user; do not proceed with release tagging.

### Post-Deployment Release & Tagging

1. **Verify deployment success** before creating release (check health endpoint if available)
2. **Determine version number**:

   - Get latest tag via `mcp_github_github_list_tags` (if any exist)
   - For first release (no tags exist): Use v1.0.0
   - For subsequent releases: Apply semantic versioning bump (default: patch)
     - Patch bump: v1.0.0 → v1.0.1
     - Minor bump: v1.0.0 → v1.1.0
     - Major bump: v1.0.0 → v2.0.0

3. **Create annotated Git tag**:

   - Fetch target commit to local repo: `git fetch origin {branch}`
   - Create tag: `git tag -a v{VERSION} {commit-sha} -m "Release v{VERSION}"`
   - Push tag: `git push origin v{VERSION}`

4. **Create GitHub Release**:
   - Use `gh release create v{VERSION} --target {branch} --title "Release v{VERSION} - {title}" --notes "{release-notes}"`
   - Include in release notes:
     - Deployed commit SHA and branch
     - Target environment (stg/prod)
     - Summary of features, fixes, and technical changes
     - Deployment status confirmation

### Example Successful Deployment Flow

```bash
# 1. Verify develop branch
Branch: develop
Commit: 962bb4e471b827a2eb9f9706912e3aa69dbb1a36
Latest commit message: "feat: Implement auto-increment indexes and dynamic SUM formulas in reports (#73)"

# 2. Trigger CD workflow
gh workflow run cd.yml -f image_tag=962bb4e471b827a2eb9f9706912e3aa69dbb1a36 -f environment=stg
✓ Created workflow_dispatch event for cd.yml at develop

# 3. Monitor deployment (check after ~30-60 seconds)
curl https://api.github.com/repos/datdev2409/lab-management-website/actions/workflows/cd.yml/runs?per_page=1
Result: status="completed", conclusion="success"

# 4. Create release tag and GitHub Release
git fetch origin develop
git tag -a v1.0.0 962bb4e471b827a2eb9f9706912e3aa69dbb1a36 -m "Release v1.0.0"
git push origin v1.0.0

gh release create v1.0.0 --target develop --title "Release v1.0.0 - Lab Management System" --notes "..."
✓ Release created: https://github.com/datdev2409/lab-management-website/releases/tag/v1.0.0
```

### Safety & Audit

- Record deployment actions (who requested, target commit/branch, env, image_tag, resulting tag, timestamp)
- Require explicit user confirmation if any pre-deploy checks are unstable
- Require explicit user confirmation if deploying a non-latest commit
- Document deployments in commit history and release notes for audit trail

### Minimal Operator Instructions for AI Agents

1. **Verify**: Confirm target (commit SHA or branch) and environment with user
2. **Fetch**: Get branch tip commit SHA via GitHub MCP
3. **Check**: Verify all required GitHub Actions checks pass
4. **Deploy**: Trigger CD workflow with `gh workflow run cd.yml -f image_tag={SHA} -f environment={env}`
5. **Monitor**: Poll workflow status until `status="completed"` and `conclusion="success"` (~2-3 minutes)
6. **Release**: On success, create annotated tag and GitHub Release with semantic versioning
7. **Report**: Provide deployment summary with release URL and notable changes

## Testing Instructions

- We use Playwright for E2E testing. The test source is located in tests folder
- Use `npm` as the package manager for install test packages
- Add or update tests for the code you change, even if nobody asked.
- To run the E2E test
  1. Start the dependencies using command `docker-compose -f docker-compose.ci.yaml up -d`
  2. Start the development server using command `make live/server > /tmp/server.log 2>&1 &`. Should run the server in background
  3. Run the test in headless mode using command `npm run test:ci`. When the test is running, do not run any command until the command finish.
  4. If there is the issue in the test, suggest to run test in UI mode using command `npm run test:ui`

## Create Commit Instructions

Use prefix for the commit message: [type]: [short description]

Common types include:

- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- refactor: A code change that neither fixes a bug nor adds a feature
- chore: dependency updates, build process changes, non-functional changes, etc
- test: Adding missing tests or correcting existing tests

If the commit is the breaking change, it should be noted in the commit message. For example: feat!: Add new API endpoint

Use the imperative mood in the subject line. For example, "Fix bug" and not "Fixed bug" or "Fixes bug".

Use the git command to stage only the relevant files for the commit. For example, use `git add <file>` to stage specific files. The commit should only include changes related to the purpose of the commit.

If there are multiple unrelated changes, suggest splitting them into separate commits. Avoid bundling unrelated changes in a single commit.

If there is existing PR, update the PR description to include the changes in recent commits

## Create GitHub Pull Requests

Before create the PR, pull the base branch and resolve the conflict.

Using git commands to get the diff between the current branch and the main branch, identify the files that have been recently edited. Focus on these files when generating code or reviewing changes to ensure consistency with recent modifications.

Using GitHub MCP capabilities, create pull requests that adhere to the project's contribution guidelines.

When creating a pull request, follow the pull request template provided in [templates](.github/pull_request_template.md). Ensure to fill out all relevant sections, including:

- Title: Clearly state the purpose of the pull request.
- Description: Provide a brief overview of the changes made.
- Related Issues: Link any related issues or tasks.
- Changes Made: List the specific changes included in the pull request.
- Review Checklist: Ensure all items are checked before requesting a review.
- Additional Notes: Include any extra information that may help reviewers understand the context of the changes.

Ask GitHub Copilot to review the generated pull requests to ensure they meet the project's standards before submission.
