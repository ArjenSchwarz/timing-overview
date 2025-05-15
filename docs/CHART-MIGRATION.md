# Chart Library Migration Plan

This document outlines the phased migration plan for switching from the deprecated `wcharczuk/go-chart` library to the recommended `vicanso/go-charts` library for pie chart generation in the Timing Overview Lambda function.

## Background

The current implementation uses the deprecated `wcharczuk/go-chart` library to generate pie charts from Timing API data. Based on the evaluation in [Charting-libraries-comparison.md](./Charting-libraries-comparison.md), `vicanso/go-charts` has been selected as the most suitable replacement due to its:

- Direct support for static PNG generation
- Native pie chart support with customization options
- Active maintenance
- Similar API structure to the current library
- Multiple themes and customization options

## Migration Phases

### Phase 1: Setup and Preparation

1. **Add the new dependency**:
   ```bash
   go get github.com/vicanso/go-charts/v2
   ```

2. **Create a test dataset**:
   - Create a JSON file with sample timing data for testing chart generation
   - This will allow for consistent comparison between the old and new implementations

3. **Set up a test environment**:
   - Create a simple test harness that can generate charts using both libraries
   - Ensure the test environment can run locally without requiring API access

### Phase 2: Direct Implementation

1. **Replace the existing chart generation function**:
   - Update `CreateProjectOverviewPieChart` in `parser/graphs.go` to use `vicanso/go-charts` instead of `wcharczuk/go-chart`
   - Implement the basic structure using `vicanso/go-charts`

2. **Implement color conversion**:
   - Create a utility function to convert hex color strings to the format required by `vicanso/go-charts`
   - Ensure project colors are correctly mapped

3. **Match existing functionality**:
   - Ensure the new implementation produces charts with the same:
     - Size dimensions (768x768)
     - Title formatting
     - Label formatting (including duration strings)
     - Color scheme

4. **Add theme support**:
   - Implement a default theme that matches the current chart appearance
   - Add support for additional themes (light, dark, grafana) as optional enhancements

### Phase 3: Testing and Validation

1. **Unit testing**:
   - Create unit tests for the updated chart generation function
   - Ensure all edge cases are covered (empty data, single item, many items)

2. **Visual validation**:
   - Generate charts using the current configuration
   - Verify the visual output meets expectations
   - Document any visual differences from the previous implementation and determine if adjustments are needed

3. **Performance testing**:
   - Measure the performance of the new implementation
   - Document memory usage and execution time

4. **Lambda testing**:
   - Deploy a test version of the Lambda function using the new implementation
   - Verify that the function works correctly in the AWS environment

### Phase 4: Deployment

1. **Direct deployment**:
   - Deploy the updated Lambda function with the new chart implementation
   - Monitor for any issues or differences in output

2. **Rollback plan**:
   - Maintain a backup of the previous implementation
   - Be prepared to roll back if any critical issues are discovered

### Phase 5: Cleanup and Documentation

1. **Remove old dependency**:
   - Remove `wcharczuk/go-chart` from go.mod
   - Run `go mod tidy` to clean up dependencies

2. **Code cleanup**:
   - Remove the old chart generation function
   - Remove any unused utility functions or variables
   - Rename the new function to match the original function name

3. **Update documentation**:
   - Update README.md with information about the new charting library
   - Document any new features or capabilities
   - Update the CHANGELOG.md with details of the migration

4. **Final review**:
   - Conduct a final code review to ensure clean implementation
   - Verify that all tests pass with the new implementation

## Detailed Implementation Steps

### Phase 2: Implementation Details

The current implementation in `parser/graphs.go` uses this pattern:

```go
func CreateProjectOverviewPieChart(configuration config.Configuration, target io.Writer) {
    // Get data from Timing API
    // Process data into chart values
    // Create and render pie chart
}
```

The updated implementation will maintain the same function signature but use the `vicanso/go-charts` API internally:

```go
func CreateProjectOverviewPieChart(configuration config.Configuration, target io.Writer) {
    // Get data from Timing API (same as before)
    // Process data into format required by go-charts
    // Create and render pie chart using go-charts
}
```

Key differences to address:

1. **Value representation**:
   - Current: Uses `chart.Value` with embedded styling
   - New: Uses separate values and styling arrays

2. **Color handling**:
   - Current: Uses `drawing.ColorFromHex`
   - New: Uses different color format/function

3. **Chart rendering**:
   - Current: Uses `chart.PNG` renderer
   - New: Uses different rendering approach

## Rollback Plan

If issues are encountered during the migration:

1. Restore the previous version of the Lambda function
2. Document the issues encountered
3. Fix issues in the new implementation
4. Retry the deployment process

## Timeline

- Phase 1: 1 day
- Phase 2: 2 days
- Phase 3: 2 days
- Phase 4: 1 day (plus monitoring period)
- Phase 5: 1 day

Total estimated time: 7 days

## Success Criteria

The migration will be considered successful when:

1. The new implementation produces visually equivalent charts
2. Performance is equal to or better than the original implementation
3. All tests pass with the new implementation
4. The old dependency has been completely removed
5. Documentation has been updated to reflect the changes
