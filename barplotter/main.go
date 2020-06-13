// Trying the creation of a new gonum/plot plotter following
// https://github.com/gonum/plot/wiki/Creating-Custom-Plotters:-A-tutorial-on-creating-custom-Plotters
package main

import (
	"fmt"
	
	"gonum.org/v1/plot/plotter"
)

// Binned bar chart
type BinnedBars struct {
	
	// Binning, ie slice of 2D array
	Binning Binning
	
	// Height of bars, ie float for every bins
	plotter.Values
}

// Trying the new plotter
func main() {

	// Test binning
	bins := NewBinning(10, 0, 10)
	fmt.Println(bins)

	
}

// 
func NewBinnedBars() BinnedBars {
	
}

// Cheap binning type - to start with something
type Binning [][2]float64

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
