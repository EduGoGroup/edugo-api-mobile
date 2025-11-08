# Implementation Plan

- [x] 1. Create bootstrap package structure and interfaces
  - Create `internal/bootstrap` directory
  - Define core interfaces in `interfaces.go` for LoggerFactory, DatabaseFactory, MessagingFactory, StorageFactory
  - Define `S3Storage` interface (reuse from existing `s3/interface.go`)
  - Define `Resources` struct to encapsulate all infrastructure resources
  - Define `BootstrapOptions` struct for configuration
  - _Requirements: 2.1, 2.2, 2.3, 4.3_

- [x] 2. Implement resource configuration system
  - [x] 2.1 Create `config.go` with resource configuration types
    - Implement `ResourceConfig` struct with Name, Optional, Enabled fields
    - Implement `DefaultResourceConfig()` function
    - _Requirements: 3.1, 3.4_

  - [x] 2.2 Implement functional options pattern
    - Create `BootstrapOption` function type
    - Implement `WithLogger()`, `WithPostgreSQL()`, `WithMongoDB()`, `WithRabbitMQ()`, `WithS3Client()` options
    - Implement `WithOptionalResource()` and `WithDisabledResource()` options
    - _Requirements: 4.1, 4.2, 3.1_

- [x] 3. Implement noop resource implementations
  - [x] 3.1 Create noop publisher
    - Create `internal/bootstrap/noop/publisher.go`
    - Implement `NoopPublisher` struct with `Publish()` and `Close()` methods
    - Add debug logging for attempted operations
    - _Requirements: 3.3_

  - [x] 3.2 Create noop storage client
    - Create `internal/bootstrap/noop/storage.go`
    - Implement `NoopS3Storage` struct with presigned URL methods
    - Add debug logging for attempted operations
    - _Requirements: 3.3_

- [x] 4. Implement resource factories
  - [x] 4.1 Create factories implementation
    - Create `factories.go` with `DefaultFactories` struct
    - Implement `CreateLogger()` wrapping `logger.NewZapLogger()`
    - Implement `CreatePostgreSQL()` wrapping `database.InitPostgreSQL()`
    - Implement `CreateMongoDB()` wrapping `database.InitMongoDB()`
    - Implement `CreatePublisher()` wrapping `rabbitmq.NewRabbitMQPublisher()`
    - Implement `CreateS3Client()` wrapping `s3.NewS3Client()`
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 5. Implement lifecycle management
  - [x] 5.1 Create lifecycle manager
    - Create `lifecycle.go` with `LifecycleManager` struct
    - Implement `Register()` method to register cleanup functions
    - Implement `Cleanup()` method to execute cleanups in LIFO order
    - Add error collection and logging for cleanup failures
    - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 6. Implement bootstrap orchestrator
  - [x] 6.1 Create bootstrapper core
    - Create `bootstrap.go` with `Bootstrapper` struct
    - Implement `New()` constructor accepting config and options
    - Apply functional options to `BootstrapOptions`
    - _Requirements: 1.1, 1.4_

  - [x] 6.2 Implement infrastructure initialization
    - Implement `InitializeInfrastructure()` method
    - Initialize logger (use injected or create new)
    - Initialize PostgreSQL (use injected or create new, handle as required resource)
    - Initialize MongoDB (use injected or create new, handle as required resource)
    - Initialize RabbitMQ (use injected or create new, handle as optional resource with noop fallback)
    - Initialize S3 (use injected or create new, handle as optional resource with noop fallback)
    - Register cleanup functions for each resource in LifecycleManager
    - Return `Resources`, cleanup function, and error
    - _Requirements: 1.1, 1.4, 3.1, 3.2, 3.3, 3.5, 4.1, 4.2, 5.1_

  - [x] 6.3 Add initialization logging
    - Log start of each resource initialization
    - Log success with relevant details (host, port, etc.)
    - Log warnings for optional resources that failed
    - Log errors for required resources that failed
    - Log total initialization time
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 7. Refactor main.go to use bootstrap system
  - [x] 7.1 Update main function
    - Replace direct resource initialization with bootstrap.New() and InitializeInfrastructure()
    - Remove direct calls to logger.NewZapLogger(), database.InitPostgreSQL(), database.InitMongoDB(), rabbitmq.NewRabbitMQPublisher(), s3.NewS3Client()
    - Add defer cleanup() call
    - Reduce main.go to under 50 lines (excluding comments and swagger annotations)
    - _Requirements: 1.1, 1.2, 1.3, 1.5_

  - [x] 7.2 Update container initialization
    - Modify container.NewContainer() to accept bootstrap.Resources instead of individual parameters
    - Update InfrastructureContainer to use Resources struct
    - Ensure backward compatibility with existing code
    - _Requirements: 1.1_

- [x] 8. Update S3Client to implement S3Storage interface
  - Verify that `s3.S3Client` implements the `S3Storage` interface
  - Add interface assertion in `s3/client.go` if needed
  - Update any type references in container to use interface
  - _Requirements: 2.3, 4.3_

- [x] 9. Add configuration support for optional resources
  - Update config structure to include optional resource flags if needed
  - Document how to configure optional resources via environment variables or config files
  - Add examples in README or config documentation
  - _Requirements: 3.4_

- [x] 10. Create integration test for bootstrap system
  - [x] 10.1 Test normal initialization
    - Create test that initializes all resources successfully
    - Verify all resources are non-nil
    - Verify cleanup function works correctly
    - _Requirements: 1.1, 5.1, 5.2_

  - [x] 10.2 Test optional resource failure
    - Create test with invalid RabbitMQ configuration
    - Mark RabbitMQ as optional
    - Verify initialization succeeds with noop publisher
    - Verify warning is logged
    - _Requirements: 3.1, 3.2, 3.3_

  - [x] 10.3 Test required resource failure
    - Create test with invalid PostgreSQL configuration
    - Verify initialization fails with descriptive error
    - Verify partial cleanup is executed
    - _Requirements: 3.5, 5.2, 5.3_

  - [x] 10.4 Test mock injection
    - Create test with all mocked resources
    - Verify mocks are used instead of real implementations
    - Verify application works with mocks
    - _Requirements: 4.1, 4.2, 4.4, 4.5_

- [x] 11. Update documentation
  - Update README with new bootstrap usage examples
  - Document how to run app with optional resources
  - Document how to inject mocks for testing
  - Add migration guide for developers
  - _Requirements: 3.4, 4.4_

- [-] 12. Commit changes to dev branch
  - Create new branch `feature/infrastructure-bootstrap-refactor` from current branch
  - Stage all uncommitted changes (config improvements)
  - Commit config improvements with descriptive message
  - Commit bootstrap implementation with descriptive message
  - Push branch for review
  - _Requirements: All_
