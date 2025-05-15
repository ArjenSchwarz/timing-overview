# Timing Overview Lambda - Modernisation Plan

## Project Description

The Timing Overview Lambda project is an AWS Lambda function that generates visual representations (pie charts) of time tracking data from the Timing app. The application retrieves time entries from the Timing API, processes them, and generates a PNG image showing the distribution of time spent across different projects.

Key components of the project include:

- **Lambda Function**: Handles API Gateway requests, retrieves time entries, and generates charts
- **Timing SDK**: Custom package for interacting with the Timing API
- **Chart Generation**: Uses go-chart library to create pie charts of time data
- **Configuration Management**: Supports both local configuration files and AWS Parameter Store for API tokens

The application can run both locally and as an AWS Lambda function, with different configuration sources depending on the execution environment.

## Requirements and Assumptions

### Modernisation Requirements

1. **Go Version Update**: Upgrade from Go 1.15 to the latest stable version (Go 1.22)
2. **Lambda Runtime Update**: Replace the deprecated `go1.x` runtime with the recommended `provided.al2023` runtime
3. **Dependency Updates**: Update all dependencies to their latest stable versions
4. **Code Modifications**: Make minimal changes to the codebase, focusing only on what's necessary for compatibility with updated dependencies

### Compatibility Requirements

1. **API Compatibility**: Maintain compatibility with the Timing API
2. **Output Compatibility**: Ensure generated charts maintain the same format and quality
3. **Configuration Compatibility**: Preserve existing configuration methods (local file and Parameter Store)

### Testing Assumptions

1. **Local Testing**: The application can be tested locally before deployment
2. **API Access**: Valid Timing API credentials are available for testing
3. **AWS Environment**: Access to AWS for testing the Lambda function is available

### Deployment Assumptions

1. **AWS SAM**: Deployment is managed using AWS SAM (Serverless Application Model)
2. **CI/CD**: No specific CI/CD pipeline requirements are considered in this plan
3. **Permissions**: Appropriate AWS permissions are available for deployment

## Phased Development Approach

### Phase 1: Environment Setup and Analysis

- [ ] **1.1** Create a git branch for modernisation work
- [ ] **1.2** Document current dependency versions and their latest available versions
- [ ] **1.3** Identify deprecated packages and functions used in the codebase
- [ ] **1.4** Set up local development environment with the latest Go version

### Phase 2: Dependency Updates

- [ ] **2.1** Update Go version in go.mod to 1.22
- [ ] **2.2** Update aws-lambda-go to the latest version
- [ ] **2.3** Update aws-sdk-go to the latest version
- [ ] **2.4** Update wcharczuk/go-chart to the latest stable version
- [ ] **2.5** Run `go mod tidy` to clean up dependencies
- [ ] **2.6** Verify that all dependencies resolve correctly

### Phase 3: Code Modifications

- [ ] **3.1** Replace deprecated `ioutil` package usage in main.go
  - [ ] **3.1.1** Replace `ioutil.ReadAll()` with `io.ReadAll()`
- [ ] **3.2** Replace deprecated `ioutil` package usage in timingsdk/tasks.go
  - [ ] **3.2.1** Replace `ioutil.ReadAll()` with `io.ReadAll()`
- [ ] **3.3** Update import statements for any renamed packages
- [ ] **3.4** Address any breaking changes in the updated dependencies
- [ ] **3.5** Update any code patterns that are deprecated in newer Go versions

### Phase 4: Lambda Runtime Update

- [ ] **4.1** Update template.yaml to use the provided.al2023 runtime
- [ ] **4.2** Configure the function to use the Go handler
- [ ] **4.3** Update any runtime-specific configurations
- [ ] **4.4** Verify SAM template validity

### Phase 5: Local Testing

- [ ] **5.1** Test local execution mode
  - [ ] **5.1.1** Verify chart generation works locally
  - [ ] **5.1.2** Confirm date parsing and formatting still works correctly
  - [ ] **5.1.3** Test with sample configuration file
- [ ] **5.2** Test Lambda execution simulation locally
  - [ ] **5.2.1** Use SAM local invoke to test Lambda function
  - [ ] **5.2.2** Verify API Gateway request handling
  - [ ] **5.2.3** Confirm Parameter Store integration works correctly

### Phase 6: Deployment and Verification

- [ ] **6.1** Build the updated Lambda package
- [ ] **6.2** Deploy to AWS using SAM
- [ ] **6.3** Verify Lambda function execution in AWS environment
- [ ] **6.4** Test API Gateway integration
- [ ] **6.5** Monitor for any errors or performance issues

### Phase 7: Documentation and Finalisation

- [ ] **7.1** Update README.md with new requirements and versions
- [ ] **7.2** Document any changes to the deployment process
- [ ] **7.3** Update any version references in documentation
- [ ] **7.4** Create pull request or merge changes to main branch
- [ ] **7.5** Tag a new release version

## Rollback Plan

In case of issues with the modernised version:

1. Revert to the previous Go version in go.mod
2. Restore the original Lambda runtime in template.yaml
3. Revert code changes related to deprecated packages
4. Rebuild and redeploy using the original configuration

## Future Considerations

- ✅ Migrated from aws-sdk-go to aws-sdk-go-v2 for improved performance and features
- Evaluate alternatives to the go-chart library if it becomes unmaintained
- Consider implementing automated tests to simplify future updates
- Explore containerisation options for more consistent local and cloud execution
