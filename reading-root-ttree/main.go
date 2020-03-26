// Exemple of reading a ROOT TTree computing some spin-correlation observables
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strings"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"github.com/golang/geo/r3"
	"github.com/rmadar/go-lorentz-vector/lv"
)

type Event struct {
	t Particles
	b Particles
	W Particles
	l Particles
	v Particles
}

type Particles struct {
	pt  []float32
	eta []float32
	phi []float32
	m   []float32
	pid []float32
}

type SpinObservables struct {
	kVec    [3]float64
	rVec    [3]float64
	nVec    [3]float64
	dphi_ll float64
	cos_km  float64
	cos_kp  float64
	cos_rm  float64
	cos_rp  float64
	cos_nm  float64
	cos_np  float64
}

func main() {

	var (
		fname   = flag.String("f", "ttbar_0j_parton.root", "path to ROOT file to analyze")
		tname   = flag.String("t", "spinCorrelation", "ROOT Tree name to analyze")
		evtmax  = flag.Int64("n", 10000, "number of events to analyze")
		verbose = flag.Bool("v", false, "verbose mode")
	)

	flag.Parse()

	eventLoop(*fname, *tname, *evtmax, *verbose)
}

// Event loop
func eventLoop(fname string, tname string, evtmax int64, verbose bool) {

	// Open the root file and get the tree
	fmt.Println("Processing the TTree", tname, "in the ROOT file", fname)
	file := openRootFile(fname)
	tree := getTtree(file, tname)

	// Create a scanner to perform the event loop on the input tree
	var (
		e Event
		rvars = getReadVariables(&e)
	)
	sc, err := rtree.NewScannerVars(tree, rvars...)
	if err != nil {
		log.Fatalf("could not create scanner: %+v", err)
	}
	defer sc.Close()
	
	// Ceate a new file, new writer to save new variables in a tree
	fnameOut := strings.ReplaceAll(fname, ".root", "_processed.root")
	fout, err := groot.Create(fnameOut)
	if err != nil {
		log.Fatalf("could not create ROOT file %q: %w", fnameOut, err)
	}
	defer fout.Close()
	var spin_var SpinObservables
	wvars := []rtree.WriteVar{
		{Name: "kvec"   , Value: &spin_var.kVec},
		{Name: "rvec"   , Value: &spin_var.rVec},
		{Name: "nvec"   , Value: &spin_var.nVec},
		{Name: "dphi_ll", Value: &spin_var.dphi_ll},
		{Name: "cosO_km", Value: &spin_var.cos_km},
		{Name: "cosO_kp", Value: &spin_var.cos_kp},
		{Name: "cosO_rm", Value: &spin_var.cos_rm},
		{Name: "cosO_rp", Value: &spin_var.cos_rp},
		{Name: "cosO_nm", Value: &spin_var.cos_nm},
		{Name: "cosO_np", Value: &spin_var.cos_np},
	}
	tout, err := rtree.NewWriter(fout, tname, wvars)
	if err != nil {
		log.Fatalf("could not create scanner: %+v", err)
	}
	defer tout.Close()

	// Actual event loop
	for sc.Next() && sc.Entry() < evtmax {

		// Entry index
		ievt := sc.Entry()

		// Load variable of the event
		err := sc.Scan()
		if err != nil {
			log.Fatalf("could not scan entry %d: %+v", ievt, err)
		}

		// Print the partonic event
		if ievt%100 == 0 && verbose {
			fmt.Println("\nEvent", ievt)
			printEvent(e)
		}

		// Re-computing spin observables
		var (
			tops = e.t
			leptons = e.l
			tplus_P4, tminus_P4 lv.FourVec
			lplus_P4, lminus_P4 lv.FourVec
		)
		// Get top/anti-top and lepton/ant-lepton four-vectors
		get4Vec := func(parts Particles, i int32) lv.FourVec {
			return lv.NewFourVecPtEtaPhiM(
				float64(parts.pt[i]),
				float64(parts.eta[i]),
				float64(parts.phi[i]),
				float64(parts.m[i])) 
		}		
		if tops.pid[0]>0 {
			tplus_P4  = get4Vec(tops, 0)
			tminus_P4 = get4Vec(tops, 1)
		} else {
			tplus_P4  = get4Vec(tops, 1)
			tminus_P4 = get4Vec(tops, 0)
		}
		if leptons.pid[0]>0 {
			lminus_P4 = get4Vec(leptons, 0) 
			lplus_P4  = get4Vec(leptons, 1) 
		} else {
			lminus_P4 = get4Vec(leptons, 1) 
			lplus_P4  = get4Vec(leptons, 0) 
		}
		// Perform the actual computation (spin basis and angles)
		cosTheta := computeSpinCosines(tplus_P4, tminus_P4, lplus_P4, lminus_P4)

		// Save the newly computed info into a TTree
		k, r, n := getSpinBasis(tplus_P4, tminus_P4)
		spin_var.kVec = [3]float64{k.X, k.Y, k.Z}
		spin_var.rVec = [3]float64{r.X, r.Y, r.Z}
		spin_var.nVec = [3]float64{n.X, n.Y, n.Z}
		spin_var.dphi_ll = lplus_P4.DeltaPhi(lminus_P4)
		spin_var.cos_km = cosTheta["k-"]
		spin_var.cos_rm = cosTheta["r-"]
		spin_var.cos_nm = cosTheta["n-"]
		spin_var.cos_kp = cosTheta["k+"]
		spin_var.cos_rp = cosTheta["r+"]
		spin_var.cos_np = cosTheta["n+"]

		_, err = tout.Write()
		if err != nil {
			log.Fatalf("could not write event %d: %+v", ievt, err)
		}
	}

	err = tout.Close()
	if err != nil {
		log.Fatalf("could not close tree-writer: %+v", err)
	}
	
	fmt.Println(" --> Event loop is done:", sc.Entry(), "events processed and stored in", fnameOut)
}

// Compute spin-related cosines
func computeSpinCosines(tplus, tminus, lplus, lminus lv.FourVec) map[string]float64 {

	// Get the proper basis
	k, r, n := getSpinBasis(tplus, tminus)

	// Get 3-vectors of lplus (lminus) in tplus (tmius) rest-frame
	lplusRF := lplus.ToRestFrameOf(tplus).Pvec
	lminusRF := lminus.ToRestFrameOf(tminus).Pvec

	// Fill the six cosines
	getCos := func(a, b r3.Vector, m float64) float64 {
		return math.Cos(a.Angle(b.Mul(m)).Radians())
	}
	cosTheta := map[string]float64{
		"k+": getCos(lplusRF, k, 1),
		"r+": getCos(lplusRF, r, 1),
		"n+": getCos(lplusRF, n, 1),
		"k-": getCos(lminusRF, k, -1),
		"r-": getCos(lminusRF, r, -1),
		"n-": getCos(lminusRF, n, -1),
	}

	return cosTheta
}

// Get spin basis
func getSpinBasis(t, tbar lv.FourVec) (k, r, n r3.Vector) {

	// Get top direction in ttbar rest frame
	ttbar := t.Add(tbar)
	top_rest := t.ToRestFrameOf(ttbar)
	k = top_rest.Pvec.Normalize()

	// Get the beam axis (Oz) and coeff to build ortho-normal basis
	beam_axis := r3.Vector{0, 0, 1}
	yval := k.Dot(beam_axis)
	ysign := yval / math.Abs(yval)
	rval := math.Sqrt(1 - yval*yval)

	// Get r axis: r = sign(y)/r * (beam -y*k)
	r = beam_axis.Add(k.Mul(-yval))
	r = r.Mul(1. / rval * ysign)

	// Get n axis: n = sign(y)/r * (beam cross k)
	n = beam_axis.Cross(k)
	n = n.Mul(1. / rval * ysign)

	return k, r, n
}


// Helper to define the event variables to load
func getReadVariables(e *Event) []rtree.ScanVar {
	return []rtree.ScanVar{
		// Top-quarks
		{Name: "t_pt", Value: &e.t.pt},
		{Name: "t_eta", Value: &e.t.eta},
		{Name: "t_phi", Value: &e.t.phi},
		{Name: "t_pid", Value: &e.t.pid},
		{Name: "t_m", Value: &e.t.m},
		// bottom-quarks
		{Name: "b_pt", Value: &e.b.pt},
		{Name: "b_eta", Value: &e.b.eta},
		{Name: "b_phi", Value: &e.b.phi},
		{Name: "b_pid", Value: &e.b.pid},
		{Name: "b_m", Value: &e.b.m},
		// W-boson
		{Name: "W_pt", Value: &e.W.pt},
		{Name: "W_eta", Value: &e.W.eta},
		{Name: "W_phi", Value: &e.W.phi},
		{Name: "W_pid", Value: &e.W.pid},
		{Name: "W_m", Value: &e.W.m},
		// Charged leptons
		{Name: "l_pt", Value: &e.l.pt},
		{Name: "l_eta", Value: &e.l.eta},
		{Name: "l_phi", Value: &e.l.phi},
		{Name: "l_pid", Value: &e.l.pid},
		{Name: "l_m", Value: &e.l.m},
		// Neutrinos
		{Name: "v_pt", Value: &e.v.pt},
		{Name: "v_eta", Value: &e.v.eta},
		{Name: "v_phi", Value: &e.v.phi},
		{Name: "v_pid", Value: &e.v.pid},
		{Name: "v_m", Value: &e.v.m},
	}
}

// Event printing
func printEvent(e Event) {
	fmt.Println(" * Top quarks")
	fmt.Println("   - pT : ", e.t.pt)
	fmt.Println("   - Eta: ", e.t.eta)
	fmt.Println("   - Phi: ", e.t.phi)
	fmt.Println("   - pid: ", e.t.pid)
	fmt.Println(" * Bottom quarks")
	fmt.Println("   - pT : ", e.b.pt)
	fmt.Println("   - Eta: ", e.b.eta)
	fmt.Println("   - Phi: ", e.b.phi)
	fmt.Println("   - pid: ", e.b.pid)
	fmt.Println(" * W-bosons")
	fmt.Println("   - pT : ", e.W.pt)
	fmt.Println("   - Eta: ", e.W.eta)
	fmt.Println("   - Phi: ", e.W.phi)
	fmt.Println("   - pid: ", e.W.pid)
	fmt.Println(" * Leptons")
	fmt.Println("   - pT : ", e.l.pt)
	fmt.Println("   - Eta: ", e.l.eta)
	fmt.Println("   - Phi: ", e.l.phi)
	fmt.Println("   - pid: ", e.l.pid)
	fmt.Println(" * Neutrinos")
	fmt.Println("   - pT : ", e.v.pt)
	fmt.Println("   - Eta: ", e.v.eta)
	fmt.Println("   - Phi: ", e.v.phi)
	fmt.Println("   - pid: ", e.v.pid)

}

// Helper to open a ROOT file
func openRootFile(fname string) *groot.File {
	f, err := groot.Open(fname)
	if err != nil {
		log.Fatalf("could not open ROOT file: %+v", err)
	}
	return f
}

// Helper to get a TTree
func getTtree(f *groot.File, tname string) rtree.Tree {
	obj, err := f.Get(tname)
	if err != nil {
		log.Fatalf("could not retrieve tree %q: %+v", tname, err)
	}
	return obj.(rtree.Tree)
}
