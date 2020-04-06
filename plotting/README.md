## Cheat sheet of go plotting

Label size and positions
```go
p := hplot.New()
p.Title.TextStyle.Font.Size = 18
p.X.Label.TextStyle.Font.Size = 18
```

Axis range
```go
p := hplot.New()
p.X.Min, p.X.Max = xmin, xmax
p.Y.Min, p.Y.Max = ymin, ymax
```

## To-do's

- [ ] add a sub-plot grid with histograms, line, points, functions
- [ ] add a sub-plot grid with different styles (axis, margin, padding, etc ...) for the same data