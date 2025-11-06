# AI Agent Instructions for Lab Management System

This document provides guidelines for AI agents working on this project, including development workflows, testing, production deployment and GitHub workflows.

## Dev Environments

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

## Create GitHub Pull Requests

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
