package main

import (

	"math"
	
	"gonum.org/v1/gonum/floats"
	
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {

	p, _ := plot.New()
	p.Title.Text   = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	var (
		n  = 100 
		x  = floats.Span(make([]float64, n), 0, 10)
		f1 = func(x float64) (y float64) { return math.Sin(x+5.) }
		f2 = func(x float64) (y float64) { return x+10. }
	)
	
	plotutil.AddLinePoints(p,
		"y = f1(x)", getFunctionPoints(x, f1),
		"y = f2(x)", getFunctionPoints(x, f2))
	
	// Save the plot to a PDF file.
	p.Save(4*vg.Inch, 4*vg.Inch, "points.pdf")
}

// randomPoints returns some random x, y points.
func getFunctionPoints(x []float64, f func(float64) float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))
	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = f(x[i])
	}
	return pts
}
