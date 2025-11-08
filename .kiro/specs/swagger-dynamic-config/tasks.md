# Implementation Plan

- [x] 1. Setup Git branch and project preparation
  - Create feature branch from dev: `feature/swagger-dynamic-config`
  - Verify swag CLI is installed locally
  - Review current Swagger configuration in codebase
  - _Requirements: 1.1, 2.1, 3.1, 4.1_

- [x] 2. Implement automatic Swagger regeneration on startup
  - [x] 2.1 Create regenerateSwagger function in cmd/main.go
    - Implement function that executes `swag init` using exec.Command
    - Add proper error handling with graceful degradation
    - Capture stdout/stderr for logging purposes
    - _Requirements: 1.1, 1.2, 1.3_
  
  - [x] 2.2 Integrate regeneration into application startup flow
    - Call regenerateSwagger before router setup in main()
    - Add logging for successful regeneration
    - Add warning logging if regeneration fails
    - Ensure application continues startup even if regeneration fails
    - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 3. Implement dynamic host detection for Swagger UI
  - [x] 3.1 Create custom Swagger UI configuration
    - Create new file internal/infrastructure/http/router/swagger.go
    - Implement SetupSwaggerUI function with custom HTML template
    - Add JavaScript code to detect window.location.host
    - Inject detected host into Swagger configuration
    - _Requirements: 2.1, 2.2, 2.4_
  
  - [x] 3.2 Update main.go Swagger annotations
    - Remove hardcoded @host annotation from main.go
    - Keep @BasePath as /v1
    - Verify all other global annotations are correct
    - _Requirements: 2.3, 4.1, 4.4_
  
  - [x] 3.3 Update router to use new Swagger UI setup
    - Replace ginSwagger.WrapHandler call in router.go
    - Use new SetupSwaggerUI function
    - Test that Swagger UI loads correctly
    - _Requirements: 2.1, 2.4_

- [x] 4. Update and validate Swagger annotations in handlers
  - [x] 4.1 Audit all handler files for missing annotations
    - Review auth_handler.go annotations
    - Review material_handler.go annotations
    - Review assessment_handler.go annotations
    - Review progress_handler.go annotations
    - Review stats_handler.go annotations
    - Review summary_handler.go annotations
    - Review health_handler.go annotations
    - _Requirements: 3.1, 3.5_
  
  - [x] 4.2 Add missing Swagger annotations to handlers
    - Add @Summary, @Description, @Tags to all endpoints
    - Add @Param annotations for all parameters
    - Add @Success and @Failure annotations for all responses
    - Add @Security annotations for protected endpoints
    - _Requirements: 3.1, 3.2, 3.3_
  
  - [x] 4.3 Update DTO structs with example tags
    - Add example tags to all request DTOs
    - Add example tags to all response DTOs
    - Ensure all required fields are marked with binding:"required"
    - _Requirements: 3.4_
  
  - [x] 4.4 Standardize tags and naming conventions
    - Ensure consistent tag names across all handlers
    - Use lowercase tags (auth, materials, assessments, etc.)
    - Group related endpoints under same tags
    - _Requirements: 4.2_

- [x] 5. Test Swagger functionality
  - [x] 5.1 Test automatic regeneration
    - Start application and verify swag init executes
    - Check logs for regeneration confirmation
    - Verify docs/ directory is updated with new timestamps
    - Test behavior when swag is not installed
    - _Requirements: 1.1, 1.2, 1.3_
  
  - [x] 5.2 Test dynamic host detection
    - Start application on default port 8080
    - Open Swagger UI and verify it loads
    - Verify API requests use correct host:port
    - Change port in config and restart
    - Verify Swagger UI adapts to new port
    - _Requirements: 2.2, 2.4, 2.5_
  
  - [x] 5.3 Test endpoint invocation from Swagger UI
    - Test /health endpoint (public)
    - Test /auth/login endpoint (public)
    - Test /materials endpoint with Bearer token (protected)
    - Verify all requests succeed with correct URLs
    - Test on different ports (8080, 3000, 9000)
    - _Requirements: 2.5, 3.2_
  
  - [x] 5.4 Validate generated Swagger documentation
    - Open swagger.json and verify all endpoints are documented
    - Verify all DTOs have proper schemas
    - Verify security definitions are correct
    - Check that no hardcoded hosts remain
    - _Requirements: 3.1, 3.3, 3.4, 3.5_

- [x] 6. Documentation and cleanup
  - [x] 6.1 Update README or documentation
    - Document the automatic Swagger regeneration feature
    - Add instructions for installing swag CLI
    - Document how to access Swagger UI
    - Add troubleshooting section for common issues
    - _Requirements: 4.1, 4.5_
  
  - [x] 6.2 Code review and cleanup
    - Remove any debug logging added during development
    - Ensure code follows project conventions
    - Verify all imports are used
    - Run go fmt on all modified files
    - _Requirements: 4.5_
  
  - [x] 6.3 Verify Git workflow compliance
    - Ensure all work is in feature branch
    - Verify no commits were made directly to dev
    - Prepare descriptive commit messages
    - Prepare PR description with testing notes
    - _Requirements: All_
