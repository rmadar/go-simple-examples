// Trying the creation of a new gonum/plot plotter following
// https://github.com/gonum/plot/wiki/Creating-Custom-Plotters:-A-tutorial-on-creating-custom-Plotters
package main

import (
	_ "fmt"
	"log"
	"image/color"

	"gonum.org/v1/gonum/floats"
	
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Binned bar chart
type BinnedBars struct {
	
	// Binning, ie slice of 2D array
	Binning Binning
	
	// Height of bars, ie float for every bins
	plotter.Values

	// Color of bars
	Color color.Color

	// Minimum Y
	Ymin float64
}

// Trying the new plotter
func main() {

	// Binned
	vals := []float64{0.4, 0.5, 0.35, 0.67, 0.80, 0.9, 1}
	bs := NewBinnedBars(vals, NewBinning(len(vals), 0, 10) )
	bs.Color = color.RGBA{R:196, B:128, A: 255}

	// Plot
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Binned Bar Chart"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(bs)

	// Save the plot
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "BinnedBars.png"); err != nil {
		panic(err)
	}
}

// NewBinnedBars creates a new binned bar plot
func NewBinnedBars(vals []float64, bins Binning) *BinnedBars {
	cpy, err := plotter.CopyValues( (plotter.Values)(vals) )
	if err != nil {
		log.Fatalf("cannot copy values")
	}
	
	return &BinnedBars{
		Binning: bins,
		Values:  cpy,
	}
}

// Plot implements the Plotter interface, drawing a rectangle
// defined by (xLow, xHigh) and Heights.
func (bs *BinnedBars) Plot(c draw.Canvas, plt *plot.Plot) {

	// Get coordinate transformation
	trX, trY := plt.Transforms(&c)

	// Loop over values
	for i, v := range bs.Values {
		ymin := trY(bs.Ymin)
		xlo := trX(bs.Binning[i][0])
		xhi := trX(bs.Binning[i][1])
		val := trY(v)
		
		var pts []vg.Point
		pts = append(pts, vg.Point{X: xlo, Y: ymin})
		pts = append(pts, vg.Point{X: xlo, Y: val})
		pts = append(pts, vg.Point{X: xhi, Y: val})
		pts = append(pts, vg.Point{X: xhi, Y: ymin})

		c.FillPolygon(bs.Color, c.ClipPolygonXY(pts))
	}
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (bs *BinnedBars) DataRange() (xmin, xmax, ymin, ymax float64) {
	catMin := bs.Ymin
	catMax := floats.Max(bs.Values)
	valMin := bs.Binning[0][0]
	valMax := bs.Binning[len(bs.Binning)-1][1]
	return valMin, valMax, catMin, catMax
}

// Simple binning type
type Binning [][2]float64

// NewBinning creates a new equal binning
func NewBinning(n int, xmin, xmax float64) Binning {
	res := make([][2]float64, n)
	dx := (xmax - xmin) / float64(n)
	for i := 0 ; i < n ; i++ {
		lo := xmin + float64(i)*dx
		hi := lo + dx
		res[i] = [2]float64{lo, hi}
	}
	return res
}
