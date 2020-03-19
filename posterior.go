package main


import (
	"fmt"
	"math"
	//"gonum.org/v1/gonum/stat"
)


func compute_mean (x []float32) float32 {
	var res, N float32 = 0.0, float32(len(x))
	for i:=0; i<len(x); i++ {
		res += x[i]/N
		fmt.Println(res)
	}
	return res	
}


func likelihood(data, model []float64) float64 {

	Nbin, Nbin2 := len(data), len(model)
	if Nbin != Nbin2 {
		fmt.Println("OOppps")
		return -1.
	}

	var LH float64 = 1.0
	for i:=0 ; i<Nbin ; i++ {
		LHbin := poisson(data[i], model[i])
		LH *= LHbin
	}
	
	return LH
}


func poisson(d, lambda float64) float64 {
	return math.Pow(lambda, d)/math.Gamma(d+1) * math.Exp(-d)
}


func create_pseudodata(model []float64) []float64 {
	pseudo_data := make([]float64, len(model))
	for ib:=0 ; ib<len(model) ; ib++ {
		pseudo_data[ib] = model[ib]
	}
	return pseudo_data
}


func main(){

	// Expectation and observation
	obs := []float64 {102, 135, 132, 125, 108};
	bkg := []float64 {100, 140, 130, 120, 110};
	sig := []float64 { 10,  35,  50,  35,   5};
	
	// Compute and print the likelihood for pseudo-data
	Ntoys := 10
	for i:=0 ; i<Ntoys ; i++ {
		pseudo_data := create_pseudodata(bkg)
		fmt.Println("LH(toys)", likelihood(pseudo_data, bkg))
	}

	// Compute and print the likelihood for observed data
	fmt.Println("LH(obs)", likelihood(obs, exp))
}

