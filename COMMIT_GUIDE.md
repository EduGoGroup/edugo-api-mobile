# Commit Guide for Swagger Dynamic Configuration Feature

This document provides suggested commit messages for the changes made in this feature branch.

## Recommended Commit Strategy

### Option 1: Single Comprehensive Commit (Recommended for Small Features)

```bash
git add .
git commit -m "feat(swagger): implement automatic regeneration and dynamic host detection

- Add automatic Swagger regeneration on application startup
- Implement dynamic host detection for Swagger UI
- Remove hardcoded @host annotations
- Add comprehensive Swagger annotations to all handlers
- Add example values to all DTOs
- Update README with Swagger documentation and troubleshooting

BREAKING CHANGE: None - backward compatible enhancement

Closes: swagger-dynamic-config spec
"
```

### Option 2: Multiple Atomic Commits (Recommended for Larger Features)

If you prefer more granular commits, use this sequence:

#### Commit 1: Core Swagger Regeneration
```bash
git add cmd/main.go
git commit -m "feat(swagger): add automatic documentation regeneration on startup

- Implement regenerateSwagger() function in cmd/main.go
- Execute swag init before HTTP server initialization
- Add graceful error handling with degradation
- Log success/failure with detailed output

Requirements: 1.1, 1.2, 1.3, 1.4
"
```

#### Commit 2: Dynamic Host Detection
```bash
git add internal/infrastructure/http/router/swagger.go internal/infrastructure/http/router/router.go
git commit -m "feat(swagger): implement dynamic host detection for Swagger UI

- Create custom Swagger UI configuration with JavaScript
- Detect window.location.host at runtime
- Remove hardcoded @host annotation from main.go
- Inject detected host into Swagger configuration

Requirements: 2.1, 2.2, 2.3, 2.4, 2.5
"
```

#### Commit 3: Handler Annotations
```bash
git add internal/infrastructure/http/handler/*.go
git commit -m "docs(swagger): add comprehensive annotations to all handlers

- Add complete Swagger annotations to auth_handler.go
- Add complete Swagger annotations to material_handler.go
- Add complete Swagger annotations to assessment_handler.go
- Add complete Swagger annotations to progress_handler.go
- Add complete Swagger annotations to stats_handler.go
- Add complete Swagger annotations to summary_handler.go
- Add complete Swagger annotations to health_handler.go
- Standardize tags and naming conventions

Requirements: 3.1, 3.2, 3.3, 3.5
"
```

#### Commit 4: DTO Examples
```bash
git add internal/application/dto/*.go
git commit -m "docs(swagger): add example values to all DTOs

- Add example tags to auth_dto.go
- Add example tags to material_dto.go
- Add example tags to stats_dto.go
- Ensure all required fields marked with binding:\"required\"

Requirements: 3.4
"
```

#### Commit 5: Documentation
```bash
git add README.md
git commit -m "docs: update README with Swagger documentation and troubleshooting

- Document automatic Swagger regeneration feature
- Add installation instructions for swag CLI
- Add comprehensive troubleshooting section
- Update development workflow documentation

Requirements: 4.1, 4.5
"
```

#### Commit 6: Cleanup (if needed)
```bash
git add .
git commit -m "chore: code cleanup and formatting

- Run go fmt on all modified files
- Remove unused imports
- Verify code follows project conventions
"
```

## Commit Message Convention

This project follows the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring without changing functionality
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates, etc.

### Scopes
- `swagger`: Swagger/OpenAPI related changes
- `config`: Configuration changes
- `api`: API endpoint changes
- `docs`: Documentation changes

### Subject
- Use imperative mood ("add" not "added" or "adds")
- Don't capitalize first letter
- No period at the end
- Maximum 50 characters

### Body (Optional)
- Explain what and why, not how
- Wrap at 72 characters
- Separate from subject with blank line

### Footer (Optional)
- Reference issues, PRs, or specs
- Note breaking changes

## Current Branch Status

```bash
# Current branch
feature/swagger-dynamic-config

# Files to commit
M  README.md
M  cmd/main.go
M  internal/application/dto/auth_dto.go
M  internal/application/dto/material_dto.go
M  internal/application/dto/stats_dto.go
M  internal/infrastructure/http/handler/assessment_handler.go
M  internal/infrastructure/http/handler/auth_handler.go
M  internal/infrastructure/http/handler/health_handler.go
M  internal/infrastructure/http/handler/material_handler.go
M  internal/infrastructure/http/handler/progress_handler.go
M  internal/infrastructure/http/handler/stats_handler.go
M  internal/infrastructure/http/handler/summary_handler.go
M  internal/infrastructure/http/router/router.go
?? .kiro/specs/swagger-dynamic-config/
?? internal/infrastructure/http/router/swagger.go
```

## Next Steps

1. Review all changes: `git diff`
2. Choose commit strategy (single or multiple commits)
3. Stage and commit files according to chosen strategy
4. Push to remote: `git push origin feature/swagger-dynamic-config`
5. Create Pull Request to `dev` branch
6. Use `PR_DESCRIPTION.md` as the PR description template

## Verification Before Commit

Run these commands to verify everything is ready:

```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run tests (if applicable)
go test ./...

# Verify application starts
go run cmd/main.go
# Should see: "documentación Swagger regenerada exitosamente"

# Check Swagger UI
# Open: http://localhost:8080/swagger/index.html
```
