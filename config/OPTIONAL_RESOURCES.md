# Optional Resources Configuration Guide

## Overview

The bootstrap system supports marking infrastructure resources as optional. When a resource is marked as optional and fails to initialize, the application will:

1. Log a warning message
2. Continue initialization with a noop (no-operation) implementation
3. Allow the application to run with reduced functionality

This is particularly useful for:
- Local development without full infrastructure
- Testing environments
- Gradual infrastructure rollout

## Supported Optional Resources

| Resource | Default | Description |
|----------|---------|-------------|
| `rabbitmq` | Optional | Message queue for async operations |
| `s3` | Optional | Object storage for file uploads |
| `postgresql` | Required | Primary database (cannot be optional) |
| `mongodb` | Required | Document database (cannot be optional) |

## Configuration Methods

### Method 1: YAML Configuration Files

Edit the appropriate config file for your environment:

**Local Development (`config-local.yaml`):**
```yaml
bootstrap:
  optional_resources:
    rabbitmq: true  # Optional - app works without it
    s3: true        # Optional - app works without it
```

**Production (`config-prod.yaml`):**
```yaml
bootstrap:
  optional_resources:
    rabbitmq: false  # Required - app fails if unavailable
    s3: false        # Required - app fails if unavailable
```

### Method 2: Environment Variables

Override configuration using environment variables:

```bash
# Make RabbitMQ required (fail if unavailable)
export BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ=false

# Make S3 required (fail if unavailable)
export BOOTSTRAP_OPTIONAL_RESOURCES_S3=false

# Run the application
go run cmd/main.go
```

**One-liner:**
```bash
BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ=false \
BOOTSTRAP_OPTIONAL_RESOURCES_S3=false \
go run cmd/main.go
```

### Method 3: Programmatic Configuration (Testing)

For integration tests, use functional options:

```go
import "github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"

// Make all resources optional for testing
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),
)

resources, cleanup, err := b.InitializeInfrastructure(ctx)
defer cleanup()
```

## Usage Examples

### Example 1: Local Development Without RabbitMQ

**Scenario:** You want to develop locally but don't have RabbitMQ running.

**Solution:**

1. Ensure `config-local.yaml` has:
```yaml
bootstrap:
  optional_resources:
    rabbitmq: true
```

2. Run the application:
```bash
APP_ENV=local go run cmd/main.go
```

3. The application will:
   - Log a warning: "RabbitMQ initialization failed, using noop implementation"
   - Continue running
   - Use a noop publisher that logs attempts but doesn't send messages

### Example 2: QA Environment Requires All Resources

**Scenario:** QA environment must have all infrastructure available.

**Solution:**

1. Configure `config-qa.yaml`:
```yaml
bootstrap:
  optional_resources:
    rabbitmq: false  # Required
    s3: false        # Required
```

2. Run the application:
```bash
APP_ENV=qa go run cmd/main.go
```

3. The application will:
   - Fail immediately if RabbitMQ is unavailable
   - Fail immediately if S3 is unavailable
   - Provide clear error messages about which resource failed

### Example 3: Temporary Override for Testing

**Scenario:** You want to test without S3 temporarily, even in an environment where it's normally required.

**Solution:**

```bash
# Override config to make S3 optional
BOOTSTRAP_OPTIONAL_RESOURCES_S3=true APP_ENV=dev go run cmd/main.go
```

### Example 4: Integration Test with Mock Resources

**Scenario:** Write integration tests without real infrastructure.

**Solution:**

```go
func TestAPIWithMocks(t *testing.T) {
    ctx := context.Background()
    cfg := testConfig()

    // Create mock resources
    mockLogger := logger.NewZapLogger("debug", "text")
    mockDB := setupTestDB(t)
    mockMongoDB := setupTestMongoDB(t)
    mockPublisher := noop.NewNoopPublisher(mockLogger)
    mockS3 := noop.NewNoopS3Storage(mockLogger)

    // Bootstrap with all mocks
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )

    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    require.NoError(t, err)
    defer cleanup()

    // Run tests with mocked infrastructure
    // ...
}
```

## Behavior Details

### When a Resource is Optional (true)

1. **Initialization Attempt:** Bootstrap tries to initialize the resource
2. **On Failure:**
   - Logs warning with error details
   - Creates noop implementation
   - Continues with other resources
3. **Runtime Behavior:**
   - Noop implementations log attempted operations
   - Operations complete without errors but don't perform actual work
   - Application remains functional with reduced capabilities

### When a Resource is Required (false)

1. **Initialization Attempt:** Bootstrap tries to initialize the resource
2. **On Failure:**
   - Logs error with full context
   - Executes cleanup for already-initialized resources
   - Returns error and stops application startup
3. **Runtime Behavior:**
   - Application never starts if resource is unavailable
   - Clear error message indicates which resource failed

## Noop Implementations

### Noop RabbitMQ Publisher

```go
// Publish logs the attempt but doesn't send messages
func (p *NoopPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
    p.logger.Debug("noop publisher: message not published (RabbitMQ not available)",
        zap.String("exchange", exchange),
        zap.String("routing_key", routingKey),
    )
    return nil
}
```

**Impact:** Events won't be published to message queues. Async operations won't trigger.

### Noop S3 Storage

```go
// GeneratePresignedUploadURL logs the attempt but returns error
func (s *NoopS3Storage) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
    s.logger.Debug("noop storage: presigned URL not generated (S3 not available)",
        zap.String("key", key),
    )
    return "", fmt.Errorf("S3 not available")
}
```

**Impact:** File uploads/downloads won't work. API endpoints requiring S3 will return errors.

## Best Practices

### Development

✅ **DO:**
- Make RabbitMQ and S3 optional in local development
- Use noop implementations for faster iteration
- Document which features require which resources

❌ **DON'T:**
- Make databases (PostgreSQL, MongoDB) optional
- Rely on noop implementations for critical features
- Forget to test with real infrastructure before deploying

### Production

✅ **DO:**
- Make all resources required in production
- Monitor resource availability
- Have alerts for initialization failures

❌ **DON'T:**
- Use optional resources in production unless intentional
- Deploy without testing resource availability
- Ignore warnings about failed optional resources

### Testing

✅ **DO:**
- Use mocks for unit tests
- Use optional resources for integration tests when appropriate
- Test both with and without optional resources

❌ **DON'T:**
- Mock everything in integration tests
- Skip testing with real infrastructure
- Assume noop implementations match real behavior

## Troubleshooting

### Problem: Application starts but features don't work

**Cause:** Optional resource failed to initialize but app continued with noop.

**Solution:**
1. Check logs for warnings about failed resources
2. Verify resource configuration (URLs, credentials)
3. Make resource required if it's critical: `bootstrap.optional_resources.rabbitmq: false`

### Problem: Application fails to start in development

**Cause:** Resource is marked as required but not available locally.

**Solution:**
1. Make resource optional in `config-local.yaml`
2. Or start the required infrastructure (Docker Compose)
3. Or use programmatic override: `BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ=true`

### Problem: Tests fail with "resource not available"

**Cause:** Test is trying to use a noop implementation that returns errors.

**Solution:**
1. Mock the resource properly in tests
2. Or use real infrastructure for integration tests
3. Or adjust test expectations for noop behavior

## Environment-Specific Defaults

| Environment | RabbitMQ | S3 | Rationale |
|-------------|----------|-----|-----------|
| `local` | Optional | Optional | Facilitate local development |
| `dev` | Optional | Optional | Allow flexible testing |
| `qa` | Required | Required | Ensure full functionality testing |
| `prod` | Required | Required | All features must work |

## Migration Guide

### Upgrading from Previous Versions

If you're upgrading from a version without optional resource configuration:

1. **No action required** - defaults match previous behavior
2. **Optional:** Add explicit configuration to your YAML files
3. **Optional:** Use environment variables for runtime overrides

### Adding New Optional Resources

To add support for a new optional resource:

1. Add field to `OptionalResourcesConfig` in `internal/config/config.go`
2. Add default in `setDefaults()` in `internal/config/loader.go`
3. Add binding in `bindEnvVars()` in `internal/config/loader.go`
4. Update bootstrap logic to check the new field
5. Create noop implementation if needed
6. Update this documentation

## Related Documentation

- [Configuration Guide](./README.md) - General configuration documentation
- [Bootstrap Design](../.kiro/specs/infrastructure-bootstrap-refactor/design.md) - Technical design details
- [Environment Variables](./README.md#variables-de-ambiente) - All available environment variables
