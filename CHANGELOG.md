# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Migrated from aws-sdk-go to aws-sdk-go-v2 for improved performance and modern capabilities
- Updated SSM Parameter Store access to use the AWS SDK v2 API
- Added context support for AWS API calls
- Created migration plan for switching from deprecated wcharczuk/go-chart to vicanso/go-charts library
- Updated chart migration plan to use direct replacement approach instead of parallel implementation

### Added
- Added vicanso/go-charts/v2 dependency for the chart library migration
- Created test dataset (test-data.json) for consistent chart generation testing
- Added test harness for comparing chart generation between old and new libraries
