# Bootstrap System Integration Tests

This document describes the integration tests for the bootstrap system.

## Overview

The integration tests validate the bootstrap system's ability to initialize infrastructure resources in various scenarios, including normal operation, optional resource failures, required resource failures, and mock injection.

## Test File

- **File**: `internal/bootstrap/bootstrap_integration_test.go`
- **Build Tag**: `integration`
- **Total Tests**: 6 integration tests (plus existing unit tests)

## Running the Tests

### Prerequisites

1. Docker must be running (for testcontainers)
2. Set the environment variable: `RUN_INTEGRATION_TESTS=true`

### Run Command

```bash
# Run all integration tests
RUN_INTEGRATION_TESTS=true go test -tags=integration -v ./internal/bootstrap

# Run specific test
RUN_INTEGRATION_TESTS=true go test -tags=integration -v ./internal/bootstrap -run TestNormalInitialization
```

## Test Coverage

### 1. TestNormalInitialization

**Requirements**: 1.1, 5.1, 5.2

**Purpose**: Validates that all resources initialize successfully with real infrastructure.

**What it tests**:
- Starts PostgreSQL, MongoDB, and RabbitMQ testcontainers
- Initializes all resources through the bootstrap system
- Verifies all resources are non-nil and functional
- Tests that PostgreSQL and MongoDB connections work
- Verifies cleanup function executes correctly
- Confirms resources are closed after cleanup

**Expected Result**: All resources initialize successfully and cleanup works properly.

---

### 2. TestOptionalResourceFailure

**Requirements**: 3.1, 3.2, 3.3

**Purpose**: Validates graceful degradation when optional resources fail.

**What it tests**:
- Starts only PostgreSQL and MongoDB (no RabbitMQ)
- Configures RabbitMQ as optional with invalid connection URL
- Verifies initialization succeeds despite RabbitMQ failure
- Confirms RabbitMQ is replaced with noop implementation
- Tests that noop publisher works without errors
- Verifies cleanup works with noop resources

**Expected Result**: Initialization succeeds with noop RabbitMQ publisher.

---

### 3. TestRequiredResourceFailure

**Requirements**: 3.5, 5.2, 5.3

**Purpose**: Validates that initialization fails when required resources fail.

**What it tests**:
- Starts only MongoDB (no PostgreSQL)
- Configures invalid PostgreSQL connection
- Verifies initialization fails with descriptive error
- Confirms error message mentions PostgreSQL
- Verifies resources and cleanup are nil on failure

**Expected Result**: Initialization fails with clear error about PostgreSQL.

---

### 4. TestRequiredResourceFailureWithPartialCleanup

**Requirements**: 3.5, 5.2, 5.3

**Purpose**: Validates partial cleanup when initialization fails mid-process.

**What it tests**:
- Starts PostgreSQL but not MongoDB
- PostgreSQL initializes successfully
- MongoDB initialization fails
- Verifies initialization fails with descriptive error
- Confirms partial cleanup is executed internally
- Verifies no panic occurs during cleanup

**Expected Result**: Initialization fails gracefully with partial cleanup.

---

### 5. TestMockInjection

**Requirements**: 4.1, 4.2, 4.4, 4.5

**Purpose**: Validates that injected mocks are used instead of real implementations.

**What it tests**:
- Creates mock resources (logger, DB, MongoDB, publisher, S3)
- Injects all mocks through bootstrap options
- Verifies initialization succeeds without real infrastructure
- Confirms all returned resources are the injected mocks
- Tests that application works with mocked resources
- Verifies cleanup works with mocks

**Expected Result**: All mocks are used, no real connections attempted.

---

### 6. TestMockInjectionWithCustomFactories

**Requirements**: 4.1, 4.2, 4.4, 4.5

**Purpose**: Validates custom factory injection for advanced testing scenarios.

**What it tests**:
- Creates custom DatabaseFactory implementation
- Tracks factory method calls
- Injects custom factory through bootstrap options
- Verifies custom factory methods are called
- Confirms factory is used instead of default implementation

**Expected Result**: Custom factory is used for resource creation.

---

## Helper Functions

### setupTestContainers
Starts PostgreSQL, MongoDB, and RabbitMQ testcontainers for full integration tests.

### setupMinimalContainers
Starts only PostgreSQL and MongoDB for testing optional resource failures.

### setupMongoOnly
Starts only MongoDB for testing PostgreSQL failure scenarios.

### setupPostgresOnly
Starts only PostgreSQL for testing MongoDB failure scenarios.

### createTestConfigWithDetails
Creates test configuration with individual connection parameters.

### createTestConfig
Creates test configuration with connection strings (for simpler tests).

### MockDatabaseFactory
Mock implementation of DatabaseFactory for testing custom factories.

## Test Execution Flow

1. **Setup**: Start required testcontainers
2. **Configure**: Create test configuration with connection details
3. **Bootstrap**: Initialize infrastructure through bootstrap system
4. **Verify**: Assert expected behavior and resource states
5. **Cleanup**: Execute cleanup function and terminate containers

## Notes

- Tests use testcontainers for real infrastructure (PostgreSQL, MongoDB, RabbitMQ)
- S3 is always configured as optional (uses noop implementation)
- Tests include 2-second delays after container startup for stability
- All tests skip if `RUN_INTEGRATION_TESTS` is not set to `true`
- Tests verify both success and failure scenarios
- Cleanup is always executed in defer blocks to prevent resource leaks

## Troubleshooting

### Tests Skip Automatically
- Ensure `RUN_INTEGRATION_TESTS=true` is set
- Check that Docker is running

### Container Startup Failures
- Verify Docker has sufficient resources
- Check that required ports are not in use
- Ensure Docker images can be pulled

### Connection Timeouts
- Increase wait times in helper functions if needed
- Check Docker network configuration
- Verify testcontainers can access Docker daemon

## Future Enhancements

- Add tests for concurrent initialization
- Add tests for resource initialization timeout scenarios
- Add tests for custom logger factory
- Add tests for custom messaging factory
- Add tests for custom storage factory
- Add performance benchmarks for initialization
