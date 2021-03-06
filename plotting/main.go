// Short example to create a plot in go
package main

import (
	"log"
	"math"
	"image/color"
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


// var defaultBlack = color.NRGBA{R: 35, G: 55, B: 57, A: 255}
var defaultBlack = color.NRGBA{R: 20, G:40 , B: 47, A: 255}

func main() {

	// Using hplot to get 1D histogram
	plotHplot_1D()

	// Using PlotUtil to get functions
	createPlotUtil()

	// Bar chart
	plotBarChart()
}


// Bar plot
func plotBarChart() {
	
	// Get some data.
	values1 := plotter.Values{0.5, 10, 20, 30}
	values2 := plotter.Values{30, 20, 10, 0.5}
	
	// Bar width
	w := 1*vg.Centimeter
	
	// Create two vertical BarCharts.
	barChart1, err := plotter.NewBarChart(values1, w) 
	if err != nil {
		log.Panic(err)
	}
	barChart2, err := plotter.NewBarChart(values2, w)
	if err != nil {
		log.Panic(err)
	}

	// Create a third vertical barChart, stack of the two.
	stack := true
	if stack {
		barChart1.StackOn(barChart2) 
	} else {
		barChart1.Offset = w
	}

	// Barchart cosmetics
	barChart1.Color = color.NRGBA{R: 200, A: 200}
	barChart1.LineStyle.Width = 0
	barChart2.Color = color.NRGBA{B: 200, A: 200}
	barChart2.LineStyle.Width = 0
	
	// Create a plot and add the bar chart.
	p := hplot.New()
	p.Add(barChart1)
	p.Add(barChart2)
	
	// Add label on the x-axis (which remove the actual axis)
	//verticalLabels := []string{"A", "B", "C", "D"}
	//p.NominalX(verticalLabels...)

	// Create a figure
	f := hplot.Figure(p)
	f.Border.Right = 15
	f.Border.Left = 5
	f.Border.Top = 10
	f.Border.Bottom = 5
	f.DPI = 150
	
	err = hplot.Save(f, 300, 100, "barChart.png")
	if err != nil {
		log.Panic(err)
	}
}

// Use hplot package
func plotHplot_1D(){

	// Get (fake) data, simulation and theoretical function
	histData, histMC, modelFunc := getData()
	
	// Make a plot and ask to compile .tex on-the-fly
	p := hplot.New()
	p.Title.Text = `\textbf{APLAS} Dummy -- $\sqrt{s}=13\,$TeV $\mathcal{L}\,=\,3\,$ab\textsuperscript{-1}`
	p.X.Label.Text = `$m_{t\bar{t}}$ [GeV]`
	p.Y.Label.Text = `$(1/\sigma) \: \mathrm{d}\sigma / \mathrm{d}m_{t\bar{t}}$`
	p.X.Label.Padding = 8
	p.Y.Label.Padding = 8
	p.X.Min, p.X.Max = -4, 6.0
	p.Y.Min, p.Y.Max =  0, 0.5
	applyPlotStyle(p) // contains all plot-related cosmetics
	//p.Add(newCustomGrid())
	
	// Create a histogram of our data values
	hData := hplot.NewH1D(histData, hplot.WithYErrBars(true))
	applyDataHistStyle(hData) // contains histogram cosmetic
	p.Add(hData)

	// Create and tune the histogram for simulation prediction
	hMC := hplot.NewH1D(histMC)
	hMC.Infos.Style = hplot.HInfoNone
	hMC.LineStyle.Width = 0
	hMC.FillColor = color.NRGBA{R: 0, G: 0, B: 100, A: 30}
	p.Add(hMC)
	
	// Tune the style of the model function (type hplot.NewFunction)
	modelFunc.Samples = 1000
	modelFunc.Color = color.NRGBA{R: 255, A: 255}
	modelFunc.Width = vg.Points(2)
	p.Add(modelFunc)
	
	// Add and tune the legend
	p.Legend.Add("Data", hData)
	p.Legend.Add("Simulation", hMC)
	p.Legend.Add("Theory", modelFunc)
	p.Legend.Top, p.Legend.Left = true, false
	p.Legend.XOffs, p.Legend.YOffs = -0.5 * vg.Inch, 0
	p.Legend.Padding = 0.1 * vg.Inch
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	// Save the plot to a TEX file
	if err := p.Save(6*vg.Inch, -1, "results/h1d_plot_wLatex.tex"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
	
	// Save the plot to a PDF file (change labels without LaTeX)
	p.Title.Text = "Gaussian Phenomena"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "PDF(x)"
	if err := p.Save(6*vg.Inch, -1, "results/h1d_plot_woLatex.pdf"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}

	// Save the plot to a PNG
	c := vgimg.NewWith(
		vgimg.UseWH(6*vg.Inch, 6*vg.Inch / math.Phi),
		vgimg.UseDPI(250),
	)
	dc := draw.New(c)
	p.Draw(dc)
	saveImg(c, "results/h1d_plot.png")
}

// Create fake data
func getData() (*hbook.H1D, *hbook.H1D, *hplot.Function) {

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
	area = 0.0
	for _, bin := range histMC.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	histMC.Scale(1 / area)

	// Function correspondong to the theory prediction
	modelFunc := hplot.NewFunction(dist.Prob)
	
	return histData, histMC, modelFunc
}

// Apply a nice plot style
func applyPlotStyle(p *hplot.Plot){

	// Specify title style
	p.Title.TextStyle.Font.Size = 18
	p.Title.TextStyle.Color = defaultBlack
	p.Title.Padding = 10

	// Specify axis ranges and padding (ie distance to (0,0) point
	p.X.Padding = 5
	p.Y.Padding = 5
	
	// Specify axis label & fontsize
	p.X.Label.TextStyle.Font.Size = 16
	p.Y.Label.TextStyle.Font.Size = 16
	p.X.Label.TextStyle.Color = defaultBlack
	p.Y.Label.TextStyle.Color = defaultBlack
	p.X.Label.Padding = 10
	p.Y.Label.Padding = 10
	// p.X.Label.Position = draw.PosRight // doesn't work well with LaTeX yet
	// p.Y.Label.Position = draw.PosTop   // doesn't work well with LaTeX yet

	// Specify axis style
	p.X.LineStyle.Width = 1.05
	p.Y.LineStyle.Width = 1.05
	p.X.LineStyle.Color = defaultBlack
	p.Y.LineStyle.Color = defaultBlack

	// Specify ticks style
	p.X.Tick.LineStyle.Width = 1.05
	p.Y.Tick.LineStyle.Width = 1.05
	p.X.Tick.LineStyle.Color = defaultBlack
	p.Y.Tick.LineStyle.Color = defaultBlack
	p.X.Tick.Label.Font.Size = 11
	p.Y.Tick.Label.Font.Size = 11
	p.X.Tick.Label.Color = defaultBlack
	p.Y.Tick.Label.Color = defaultBlack

	// Specify tick position
	p.X.Tick.Marker = hplot.Ticks{N: 10}
	p.Y.Tick.Marker = hplot.Ticks{N: 10}

	// Specify text style of the legend
	p.Legend.TextStyle.Font.Size = 14
	p.Legend.TextStyle.Color = defaultBlack
}

func applyDataHistStyle(hData *hplot.H1D){

	// Remove basic stat info
	hData.Infos.Style = hplot.HInfoNone
	
	// No line
	hData.LineStyle.Width = 0

	// Y error bars
	hData.YErrs.LineStyle.Color = defaultBlack
	hData.YErrs.LineStyle.Width = 2
	hData.YErrs.CapWidth = 6

	// Dots as marker
	hData.GlyphStyle = draw.GlyphStyle{
		Shape:  draw.CircleGlyph{},
		Color:  defaultBlack,
		Radius: vg.Points(2.5)}
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
	err = p.Save(4*vg.Inch, 4*vg.Inch, "results/points.pdf")
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

