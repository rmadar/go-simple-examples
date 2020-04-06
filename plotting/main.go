// Short example to create a plot in go
package main

import (
	"log"
	"math"
	"image/color"
	"fmt"
	
	"gonum.org/v1/gonum/stat/distuv"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/floats"
	
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	_ "gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg"

	
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

func main() {

	plotHplot_1D()
	createPlotUtil()

}

// Use hplot package
func plotHplot_1D(){

	const npoints = 500

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title and axis range 
	p := hplot.New()
	p.Title.TextStyle.Font.Size = 18
	p.X.Label.TextStyle.Font.Size = 14
	p.Y.Label.TextStyle.Font.Size = 14
	p.Title.Text = "Gaussian PDF" // "Histogram $\\exp(\\frac{x^2}{\\sigma^2})$"
	p.X.Label.Text = "X"          // "$\\sqrt{X}$"
	p.Y.Label.Text = "Y"
	p.X.Min, p.X.Max = -4, 7
	p.Y.Min, p.Y.Max =  0, 0.6
	
	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist, hplot.WithYErrBars(true))
	h.Infos.Style = hplot.HInfoNone
	h.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 0}
	h.FillColor = color.RGBA{R: 90, G: 90, B: 250, A: 80}
	h.YErrs.LineStyle.Color = color.RGBA{R: 90, G: 90, B: 250, A: 255}
	h.YErrs.LineStyle.Width = +0.02 * vg.Inch
	// h.Shape = draw.CircleGlyph{} --> doesn't work	
	p.Add(h)
	
	// Add the normal distribution function with hplot.NewFunction 
	norm1 := hplot.NewFunction(dist.Prob)
	norm1.Color = color.RGBA{R: 255, A: 180} // {R: 255, A: 100} doesn't compile  with LaTeX
	norm1.Width = vg.Points(2)
	p.Add(norm1)

	// Add the normal distribution function with Function structure directly
	norm2 := &hplot.Function{F:dist.Prob, Samples: 10}
	norm2.Color = color.RGBA{G: 250, A: 180} // {R: 255, A: 100} doesn't compile  with LaTeX
	norm2.Width = vg.Points(2)
	p.Add(norm2)
	
	// Add and manipulate a legend
	p.Legend.Add(fmt.Sprintf("Sampling (%v toys)", npoints), h)
	p.Legend.Add("PDF (n=30, default)", norm1)
	p.Legend.Add("PDF (n=10, user)", norm2)
	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.YOffs = -0.25 * vg.Inch
	p.Legend.XOffs = -0.5 * vg.Inch
	p.Legend.Padding = 0.1 * vg.Inch
	p.Legend.ThumbnailWidth = 0.3 * vg.Inch
	p.Legend.TextStyle.Font.Size = 12
	
	// draw a grid if we want
	// p.Add(hplot.NewGrid())

	// Save the plot to a PDF file.
	if err := p.Save(6*vg.Inch, -1, "h1d_plot.pdf"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// Use plotUtil package
func createPlotUtil(){
	
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
		f2 = func(x float64) (y float64) { return x/2. + 2. }
	)

	// Use the plotutil package
	plotutil.AddLinePoints(p,
		"y = f1(x)", getPoints(x, f1),
		"y = f2(x)", getPoints(x, f2),
	)

	cos := plotter.NewFunction(math.Cos)
	cos.LineStyle.Color = plotutil.SoftColors[2]
	p.Add(cos)
	p.Legend.Add("cos(x)", cos)
	p.Legend.Top = true
	p.Legend.Left = true

	// Save the plot to a PDF file.
	err = p.Save(4*vg.Inch, 4*vg.Inch, "points.pdf")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

// get plotter.YXs where X=xs and ys=f(Xs)
func getPoints(xs []float64, f func(float64) float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i, x := range xs {
		pts[i].X = x
		pts[i].Y = f(x)
	}
	return pts
}
