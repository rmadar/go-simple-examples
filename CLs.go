package main


import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/gonum/stat"
)


func Filter(x []float64, condition func(float64) bool) []float64 {
	res := make([]float64, 0)
	for _, val := range x {
		if condition(val) { res = append(res, val) }
	}
	return res
}


func model_prediction(bkg, sig []float64, mu float64) []float64 {
	prediction := make([]float64, len(bkg))
	for i:=0 ; i<len(bkg) ; i++ {
		prediction[i] = bkg[i] + mu*sig[i]
	}
	return prediction
}


func likelihood(data, model []float64) float64 {
	var LH float64 = 1.0
	for i, v := range data {
		LH *= distuv.Poisson{Lambda: model[i]}.Prob(v)
	}
	return LH
}


func NLLR(data, model1, model2 []float64) float64 {
	L_hyp1 := likelihood(data, model1)
	L_hyp2 := likelihood(data, model2)
	return -2*math.Log(L_hyp1/L_hyp2)
}


func create_pseudodata(model []float64) []float64 {
	pseudo_data := make([]float64, len(model))
	for ib:=0 ; ib<len(model) ; ib++ {
		pseudo_data[ib] = distuv.Poisson{Lambda: model[ib]}.Rand()
	}
	return pseudo_data
}


func main() {
	
	// Expectation and observation
	obs := []float64 {102, 135, 132, 125, 108}
	bkg := []float64 {100, 140, 130, 120, 110}
	sig := []float64 { 10,  35,  50,  35,   5};

	// Get B-only and S+B expectations
	model_SB    := model_prediction(bkg, sig, 1.0)
	model_Bonly := model_prediction(bkg, sig, 0.0)
	
	// Compute and print the likelihood for pseudo-data
	Ntoys := 200000
	var nllr_sb, nllr_b []float64;
	for i:=0 ; i<Ntoys ; i++ {
		if i%100000==0 { fmt.Println("Toy n", i) }
		nllr_sb = append(nllr_sb, NLLR(create_pseudodata(model_SB), model_SB, model_Bonly))
		nllr_b  = append(nllr_b , NLLR(create_pseudodata(model_Bonly), model_SB, model_Bonly))
	}

	// Print some results
	fmt.Println("Mean[NLLR(S+B|B]  =", stat.Mean(nllr_sb, nil))
	fmt.Println("Mean[NLLR(B|B)]   =", stat.Mean(nllr_b , nil))
	fmt.Println("Mean[NLLR(obs|B)] =", NLLR(obs, model_SB, model_Bonly))

	// Compute & print CLs+b and CLs
	Nexp := len(Filter(nllr_b, func(x float64) bool {return x>32} ) )
	CLs := float32(Nexp) / float32(Ntoys)
	fmt.Println(CLs)
}

