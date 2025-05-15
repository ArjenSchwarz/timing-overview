# Go Charting Libraries Comparison for Pie Charts

Based on my research using the Perplexity MCP, I've identified several alternatives to the deprecated `wcharczuk/go-chart` library that support pie charts. Here's a detailed comparison focusing on the most suitable options for your project:

## Top Recommendations

### 1. **vicanso/go-charts**
This appears to be the most direct and suitable replacement for your needs.

**Strengths:**
- Direct support for static PNG generation (critical for your Lambda function)
- Native pie chart support with customization options
- Actively maintained (last commit in 2025)
- Built as an improved fork of the original go-chart library
- Multiple themes (light, dark, grafana, ant)
- Apache ECharts-compatible options
- Straightforward migration path from wcharczuk/go-chart

**Example code:**
```go
values := []float64{1048, 735, 580, 484, 300}
p, err := charts.PieRender(
  values,
  charts.TitleOptionFunc(charts.TitleOption{
    Text: "Project Overview",
    Subtext: "Time Spent",
    Left: charts.PositionCenter,
  }),
  charts.PaddingOptionFunc(charts.Box{
    Top: 20,
    Right: 20,
    Bottom: 20,
    Left: 20,
  }),
)
// Write to file
f, _ := os.Create("output.png")
p.Render(f)
```

### 2. **gonum/plot with benoitmasson/plotters/piechart**

**Strengths:**
- Scientific plotting library with strong mathematical foundations
- Native PNG generation
- Pie charts available through the extension package
- Good for statistical and scientific visualizations

**Limitations:**
- Requires an additional package for pie charts
- More complex API than vicanso/go-charts
- Less focused on business/application charts

### 3. **go-echarts**

**Strengths:**
- Rich interactive features
- Extensive customization options
- Large community and good documentation
- Based on Apache ECharts

**Limitations:**
- Primarily generates HTML/JS output
- Would require additional tools (like headless browsers) to generate static PNG images
- Migration would be more complex

## Migration Considerations

For your Lambda function that generates pie charts as PNG images, **vicanso/go-charts** offers the smoothest migration path from wcharczuk/go-chart:

1. **API Similarity**: Similar structure but with improved options
2. **Output Format**: Direct PNG generation without additional dependencies
3. **Maintenance**: Active development and community support
4. **Features**: Additional themes and customization options beyond the original library

## Not Recommended

- **junglerider/go-charts**: Appears to be an inactive fork or mirror with unclear maintenance status
- **fogleman/gg**: While powerful for 2D graphics, it would require manual implementation of pie chart functionality

## Summary Comparison Table

| Library | PNG Support | Pie Charts | Maintenance | Migration Difficulty | Additional Features |
|---------|------------|------------|-------------|----------------------|---------------------|
| vicanso/go-charts | Native | Built-in | Active | Low | Themes, ECharts compatibility |
| gonum/plot + extension | Native | Via extension | Active | Medium | Scientific plotting capabilities |
| go-echarts | Via tools | Built-in | Active | High | Interactive, web-friendly |
| wcharczuk/go-chart (current) | Native | Built-in | Deprecated | N/A | Limited |

Would you like me to provide more specific details about any of these libraries, or would you prefer to move forward with implementing a replacement using vicanso/go-charts? When you're ready to implement the changes, please toggle to Act mode.