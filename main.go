package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {

	// Expectation and observation
	obs := []float64{102, 135, 132, 125, 108}
	bkg := []float64{100, 140, 130, 120, 110}
	sig := []float64{0, 5, 20, 15, 2}

	// Run the CLs computation for these background, signal and observations
	POI, CLs_exp, CLs_obs := computeCLsVsPOI(bkg, sig, obs)

	// Print the results
	for i, mu := range POI {
		fmt.Println("\nmu =", mu)
		fmt.Println("  -> CLs[exp] =", CLs_exp[i])
		fmt.Println("  -> CLs[obs] =", CLs_obs[i])
	}
}

func computeCLsVsPOI(bkg, sig, obs []float64) (POI, CLs_exp, CLs_obs []float64) {

	// Number of pseudo-experiment per mu value
	Ntoys := 100000

	// Get B-only expectaction and associated toys
	model_Bonly := modelPrediction(bkg, sig, 0.0)
	pseudodata_Bonly := make([][]float64, 0)
	for i := 0; i < Ntoys; i++ {
		pseudodata_Bonly = append(pseudodata_Bonly, createPseudodata(model_Bonly))
	}

	// Start to loop over mu values
	nPOI := 20
	POI_start, POI_end := 0.0, 2.0
	for iPOI := 0; iPOI < nPOI; iPOI++ {

		// Get S+B expectations
		mu := POI_start + (POI_end-POI_start)/float64(nPOI)*float64(iPOI)
		POI = append(POI, mu)
		model_SB := modelPrediction(bkg, sig, mu)

		// Get observed nllr for this assumed POI value
		nllr_obs := NLLR(obs, model_SB, model_Bonly)

		// Draw some toys to get PDF(nllr|S+B) and PDF(nllr|B)
		var nllr_sb, nllr_b []float64
		for i := 0; i < Ntoys; i++ {
			nllr_sb = append(nllr_sb, NLLR(createPseudodata(model_SB), model_SB, model_Bonly))
			nllr_b = append(nllr_b, NLLR(pseudodata_Bonly[i], model_SB, model_Bonly))
		}
		CLs_exp = append(CLs_exp, computeCLs(nllr_sb, nllr_b, stat.Mean(nllr_b, nil)))
		CLs_obs = append(CLs_obs, computeCLs(nllr_sb, nllr_b, nllr_obs))
	}

	return POI, CLs_exp, CLs_obs
}

func modelPrediction(bkg, sig []float64, mu float64) []float64 {
	prediction := make([]float64, len(bkg))
	for i := 0; i < len(bkg); i++ {
		prediction[i] = bkg[i] + mu*sig[i]
	}
	return prediction
}

func NLLR(data, model1, model2 []float64) float64 {
	L_hyp1 := likelihood(data, model1)
	L_hyp2 := likelihood(data, model2)
	return -2 * math.Log(L_hyp1/L_hyp2)
}

func likelihood(data, model []float64) float64 {
	var LH float64 = 1.0
	for i, v := range data {
		LH *= distuv.Poisson{Lambda: model[i]}.Prob(v)
	}
	return LH
}

func createPseudodata(model []float64) []float64 {
	pseudo_data := make([]float64, len(model))
	for ib := 0; ib < len(model); ib++ {
		pseudo_data[ib] = distuv.Poisson{Lambda: model[ib]}.Rand()
	}
	return pseudo_data
}

func computeCLs(nllr_sb, nllr_b []float64, ref float64) float64 {
	condition := func(x float64) bool { return x >= ref }
	Nsb := len(filterSlice(nllr_sb, condition))
	Nb := len(filterSlice(nllr_b, condition))
	CLsb := float64(Nsb) / float64(len(nllr_sb))
	CLb := float64(Nb) / float64(len(nllr_b))
	return CLsb / CLb
}

func filterSlice(x []float64, condition func(float64) bool) []float64 {
	res := make([]float64, 0)
	for _, val := range x {
		if condition(val) {
			res = append(res, val)
		}
	}
	return res
}
