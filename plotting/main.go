// Short example to create a plot in go
package main

import (
	"log"
	"math"
	"image/color"
	"fmt"
	"os"
	
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


var defaultBlack color.NRGBA = color.NRGBA{R: 40, G: 40, B: 40, A: 240}

func main() {

	// Using hplot to get 1D histogram
	plotHplot_1D()

	// Using PlotUtil to get functions
	createPlotUtil()
}

// Use hplot package
func plotHplot_1D(){

	const npoints = 500

	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard normal distribution
	histData := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		histData.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range histData.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	histData.Scale(1 / area)

	// Draw some random values from the standard normal distribution
	histMC := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		histMC.Fill(v, 1)
	}

	// normalize histogram
	area = 0.0
	for _, bin := range histMC.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	histMC.Scale(1 / area)
	
	// Make a plot and set its title and axis range 
	p := hplot.New()

	// draw a grid which is produced by a function
	// p.Add(newCustomGrid())
	
	// Plot borders --> doesn't work with LaTeX/PNG with vgimg yet
	p.Border.Right = 15
	p.Border.Left = 5
	p.Border.Top = 10
	p.Border.Bottom = 5

	// Specify titles
	p.Title.Text = "Gaussian Phenomena"
	//p.Title.Text = "\\textbf{APLAS} Dummy -- $\\sqrt{s}=13\\,$TeV $\\mathcal{L}\\,=\\,3\\,$ab\\textsuperscript{-1}"
	
	// Specify title style
	p.Title.TextStyle.Font.Size = 18
	p.Title.TextStyle.Color = defaultBlack
	p.Title.Padding = 10

	// Specify axis ranges and padding (ie distance to (0,0) point
	p.X.Min, p.X.Max = -4, 6
	p.Y.Min, p.Y.Max =  0, 0.5
	p.X.Padding = 5
	p.Y.Padding = 5
	
	// Specify axis label & fontsize
	//p.X.Label.Text = "$m_{t\\bar{t}}$ [GeV]"
	//p.Y.Label.Text = "$(1/\\sigma) \\: \\mathrm{d}\\sigma / \\mathrm{d}m_{t\\bar{t}}$"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.X.Label.TextStyle.Font.Size = 18
	p.Y.Label.TextStyle.Font.Size = 18
	p.X.Label.TextStyle.Color = defaultBlack
	p.Y.Label.TextStyle.Color = defaultBlack

	// Specify axis line, ticks values, label style & size
	p.X.LineStyle.Width = 1.1
	p.Y.LineStyle.Width = 1.1
	p.X.LineStyle.Color = defaultBlack
	p.Y.LineStyle.Color = defaultBlack
	p.X.Tick.LineStyle.Width = 1.1
	p.Y.Tick.LineStyle.Width = 1.1
	p.X.Tick.LineStyle.Color = defaultBlack
	p.Y.Tick.LineStyle.Color = defaultBlack
	p.X.Tick.Marker = customTicks(p.X.Min, p.X.Max)
	p.Y.Tick.Marker = customTicks(p.Y.Min, p.Y.Max)
	p.X.Tick.Label.Font.Size = 14
	p.Y.Tick.Label.Font.Size = 14
	p.X.Tick.Label.Color = defaultBlack
	p.Y.Tick.Label.Color = defaultBlack
	
	// Create a histogram of our values drawn from the standard normal.
	h := hplot.NewH1D(histData, hplot.WithYErrBars(true))

	// Info of histogram 
	h.Infos.Style = hplot.HInfoNone

	// Line color / style
	h.LineStyle.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 30}
	h.LineStyle.Width = 0

	// Error bar style
	h.YErrs.LineStyle.Color = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
	h.YErrs.LineStyle.Width = 1.5

	// Histo marker style / color
	h.GlyphStyle = draw.GlyphStyle{
		Shape:  draw.CircleGlyph{},
		Color:  color.NRGBA{R: 10, G: 10, B: 10, A: 255},
		Radius: vg.Points(3)}
	p.Add(h)

	
	hMC := hplot.NewH1D(histMC)
	hMC.Infos.Style = hplot.HInfoNone
	hMC.LineStyle.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 30}
	hMC.LineStyle.Width = 0
	hMC.FillColor = color.NRGBA{R: 0, G: 0, B: 100, A: 30}
	p.Add(hMC)
	
	// Add the normal distribution function with hplot.NewFunction 
	norm := hplot.NewFunction(dist.Prob)
	norm.Samples = 1000
	norm.Color = color.NRGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)
	
	// Add and manipulate a legend
	p.Legend.Add("Data", h)
	p.Legend.Add("Simulation", hMC)
	p.Legend.Add("Theory", norm)
	p.Legend.Top = true
	p.Legend.Left = false
	p.Legend.YOffs = 0
	p.Legend.XOffs = -0.5 * vg.Inch
	p.Legend.Padding = 0.1 * vg.Inch
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch
	p.Legend.TextStyle.Font.Size = 15
	p.Legend.TextStyle.Color = defaultBlack

	// Save the plot to a TEX file
	if err := p.Save(6*vg.Inch, -1, "h1d_plot_wLatex.tex"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
	
	// Save the plot to a PDF file
	if err := p.Save(6*vg.Inch, -1, "h1d_plot_woLatex.pdf"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}

	// Save the plot to a PNG (resolution doesnt work resolution - most likely not used properly)
	c := vgimg.NewWith(
		vgimg.UseWH(12*vg.Centimeter, 9*vg.Centimeter),
		vgimg.UseDPI(200),
	)
	dc := draw.New(c)
	p.Draw(dc)
	saveImg(c, "h1d_plot.png")
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


func saveImg(c *vgimg.Canvas, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalf("could not create output image file: %+v", err)
	}
	defer f.Close()

	cpng := vgimg.PngCanvas{Canvas: c}
	_, err = cpng.WriteTo(f)
	if err != nil {
		log.Fatalf("could not encode image to PNG: %+v", err)
	}
	
	err = f.Close()
	if err != nil {
		log.Fatalf("could not close output image file: %+v", err)
	}
}


// Functions stolen from David Calvet:
// https://gitlab.cern.ch/atlas-clermont/tile/front-end/fengo/-/blob/master/PlotHelpers.go
func newCustomGrid() *plotter.Grid {
	gr := plotter.NewGrid()
	gr.Vertical.Color = color.Gray{220}
	gr.Vertical.Dashes = []vg.Length{vg.Points(3), vg.Points(5)}
	gr.Horizontal.Color = color.Gray{220}
	gr.Horizontal.Dashes = []vg.Length{vg.Points(3), vg.Points(5)}
	return gr
}

func customTicks(ymin, ymax float64) plot.ConstantTicks {

	ticks := []plot.Tick{}
	// computing order of range (position of least significant digit)
	yorder := int(math.Log10(ymax-ymin)+0.5) - 1
	format := ".0f"
	if yorder < 1 {
		format = fmt.Sprintf(".%df", -yorder)
	}
	// stepping is a power of 10 with integer exponent (yorder)
	ystep := math.Pow10(yorder)
	// tuning step
	if (ymax-ymin)/ystep > 20 {
		ystep *= 5
	}
	// first big tick is rounded to the correct significant digit
	yoffset := float64(int(ymin/ystep)) * ystep
	// creating big ticks
	for y := yoffset; y <= ymax; y += ystep {
		label := fmt.Sprintf("%"+format, y)
		ticks = append(ticks, plot.Tick{y, label})
	}
	// 5 small ticks for each big tick
	ysub := ystep / 5
	for y := yoffset - ysub; y >= ymin; y -= ysub {
		ticks = append(ticks, plot.Tick{y, ""})
	}
	for y := yoffset + ysub; y <= ymax; y += ysub {
		ticks = append(ticks, plot.Tick{y, ""})
	}
	return plot.ConstantTicks(ticks)
}
