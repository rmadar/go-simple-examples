package main

import (
	"log"
	"math"

	"gonum.org/v1/gonum/floats"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {

	// Create a plot object
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %+v", err)
	}
	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Define variable to be plotted
	var (
		n  = 100
		x  = floats.Span(make([]float64, n), 0, 10)
		f1 = func(x float64) (y float64) { return math.Sin(x + 5.) }
		f2 = func(x float64) (y float64) { return x + 10. }
	)

	// Use the plotutil package
	plotutil.AddLinePoints(p,
		"y = f1(x)", getPoints(x, f1),
		"y = f2(x)", getPoints(x, f2))
	
	// Save the plot to a PDF file.
	err = p.Save(4*vg.Inch, 4*vg.Inch, "points.pdf")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

func getPoints(xs []float64, f func(float64) float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i, x := range xs {
		pts[i].X = x
		pts[i].Y = f(x)
	}
	return pts
}
