## Cheat sheet of go plotting

Structure of the plotting tools (WIP):
```go
Figure
 |- canvas
    |- plot
       |- histogram
       |- functions
       |- legends
       |- axis
       |...
```


Label size and positions
```go
p := hplot.New()
p.Title.TextStyle.Font.Size = 18
p.X.Label.TextStyle.Font.Size = 18
p.X.Label.XAlign = draw.XRight // this doesn't put the label to the right (?)
```

Axis range
```go
p := hplot.New()
p.X.Min, p.X.Max = xmin, xmax
p.Y.Min, p.Y.Max = ymin, ymax
```

Adding a Legend
```go
p := hplot.New()
p.Legend.Add(obj, "name")	
p.Legend.Top, p.Legend.Left = true, false // position
p.Legend.YOffs = -0.25 * vg.Inch          // offset wrt to position
p.Legend.XOffs = -0.5 * vg.Inch           // offset wrt to position
p.Legend.Padding = 0.1 * vg.Inch          // padding of the legend
p.Legend.ThumbnailWidth = 0.3 * vg.Inch   // width of the legend
```

## To-do's

- [ ] add a sub-plot grid with histograms, line, points, functions
- [ ] add a sub-plot grid with different styles (axis, margin, padding, etc ...) for the same data