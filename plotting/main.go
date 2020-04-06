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
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"

	
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

func main() {

	// Using hplot to get 1D histogram
	plotHplot_1D()

	// Using PlotUtil to get functions
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

	// Specify titles
	p.Title.Text = "Gaussian PDF" // "Histogram $\\exp(\\frac{x^2}{\\sigma^2})$"
	p.X.Label.Text = "X"          // "$\\sqrt{X}$"
	p.Y.Label.Text = "Y"

	// Specify title styles
	p.Title.TextStyle.Font.Size = 18
	p.Title.Padding = 0
	p.X.Label.TextStyle.Font.Size = 14
	p.X.Label.XAlign = draw.XLeft
	p.Y.Label.TextStyle.Font.Size = 14
	
	// Specify axis ranges and padding
	p.X.Min, p.X.Max = -4, 7
	p.Y.Min, p.Y.Max =  0, 0.6
	p.X.Padding = 5
	p.Y.Padding = 5
		
	// Create a histogram of our values drawn from the standard normal.
	h := hplot.NewH1D(hist, hplot.WithYErrBars(true))
	h.Infos.Style = hplot.HInfoNone
	h.LineStyle.Color = color.NRGBA{R: 10, G: 10, B: 200, A: 255}
	h.LineStyle.Width = 2
	h.FillColor = color.NRGBA{R: 90, G: 90, B: 200, A: 100}
	h.YErrs.LineStyle.Color = color.NRGBA{R: 10, G: 10, B: 200, A: 255}
	h.YErrs.LineStyle.Width = 2
	p.Add(h)
	
	// Add the normal distribution function with hplot.NewFunction 
	norm1 := hplot.NewFunction(dist.Prob)
	norm1.Color = color.NRGBA{R: 255, A: 180}
	norm1.Width = vg.Points(2)
	p.Add(norm1)

	// Add the normal distribution function with Function structure directly
	norm2 := &hplot.Function{F:dist.Prob, Samples: 10}
	norm2.Color = color.NRGBA{G: 250, A: 180}
	norm2.Width = vg.Points(2)
	p.Add(norm2)
	
	// Add and manipulate a legend
	p.Legend.Add(fmt.Sprintf("Sampling (%v toys)", npoints), h)
	p.Legend.Add("PDF (n=30, default)", norm1)
	p.Legend.Add("PDF (n=10, user)", norm2)
	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.YOffs = -0.5 * vg.Inch
	p.Legend.XOffs = -0.5 * vg.Inch
	p.Legend.Padding = 0.1 * vg.Inch
	p.Legend.ThumbnailWidth = 0.3 * vg.Inch
	p.Legend.TextStyle.Font.Size = 12
	
	// draw a grid if we want
	// p.Add(hplot.NewGrid())

	// Save the plot to a PDF file
	if err := p.Save(6*vg.Inch, -1, "h1d_plot.pdf"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}

	// Save the plot to a PNG (resolution doesnt work resolution - most likely not used properly)
	c := vgimg.NewWith(
		vgimg.UseWH(10*vg.Centimeter, 12*vg.Centimeter),
		vgimg.UseDPI(200),
	)
	dc := draw.New(c)
	p.Draw(dc)
	if err := p.Save(6*vg.Inch, -1, "h1d_plot.png"); err != nil {
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
