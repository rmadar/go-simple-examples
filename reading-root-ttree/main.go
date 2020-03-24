// Exemple of reading a ROOT TTree computing some spin-correlation observables
package main

import (
	"fmt"
	"log"
	"math"
	
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"github.com/rmadar/go-lorentz-vector/lv"
	"github.com/golang/geo/r3"
)


type PartonEvent struct {
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
	kVec []float32 
	rVec []float32 
	nVec []float32 
	dphi_ll []float32
	cos_km []float32
	cos_kp []float32
	cos_rm []float32
	cos_rp []float32
	cos_nm []float32
	cos_np []float32
}

const mtop float64 = 173.

func main() {

	var (
		fname  = "ttbar_0j_parton.root"
		tname  = "spinCorrelation"
		evtmax = int64(10)
	)

	eventLoop(fname, tname, evtmax)
}



// Event loop
func eventLoop(fname string, tname string, evtmax int64) {

	// Open the root file and get the tree
	fmt.Println("Opening TTree", tname, "in ROOT file", fname)
	file := openRootFile(fname)
	tree := getTtree(file, tname)
	
	// Event and variables to load
	var (
		e_partons PartonEvent
		e_spin_obs SpinObservables
		vars = make([]rtree.ScanVar, 0)
	)
	vars = append(vars, getPartonVariables(&e_partons)...)
	vars = append(vars, getAngularVariables(&e_spin_obs)...)
	
	// Create a scanner to perform the event loop
	sc, err := rtree.NewScannerVars(tree, vars...)
	if err != nil {
		log.Fatalf("could not create scanner: %+v", err)
	}
	defer sc.Close()

	// Actual event loop
	for sc.Next() && sc.Entry() < evtmax {

		// Entry index
		ievt := sc.Entry()

		// Load variable of the event
		err := sc.Scan()
		if err != nil {
			log.Fatalf("could not scan entry %d: %+v", ievt, err)
		}

		// Getting slice of particles
		tops := e_partons.t
		leptons := e_partons.l
				
		// Re-computing spin observables
		var (
			loadVec   = lv.NewFourVecPtEtaPhiM
			tplus_P4  = loadVec(float64(tops.pt[0]), float64(tops.eta[0]), float64(tops.phi[0]), mtop)
			tminus_P4 = loadVec(float64(tops.pt[1]), float64(tops.eta[1]), float64(tops.phi[1]), mtop)
			lminus_P4 = loadVec(float64(leptons.pt[0]), float64(leptons.eta[0]), float64(leptons.phi[0]), 0.0)
			lplus_P4  = loadVec(float64(leptons.pt[1]), float64(leptons.eta[1]), float64(leptons.phi[1]), 0.0)
		)
		cosTheta := computeSpinCosines(tplus_P4, tminus_P4, lminus_P4, lplus_P4)

		// Compare with stored variables
		fmt.Println("\nComparing stored and recomputed spin variables:")
		fmt.Println("km: ", float64(e_spin_obs.cos_km[0]) - cosTheta["km"])
		fmt.Println("rm: ", float64(e_spin_obs.cos_rm[0]) - cosTheta["rm"])
		fmt.Println("nm: ", float64(e_spin_obs.cos_nm[0]) - cosTheta["nm"])
		fmt.Println("kp: ", float64(e_spin_obs.cos_kp[0]) - cosTheta["kp"])
		fmt.Println("rp: ", float64(e_spin_obs.cos_rp[0]) - cosTheta["rp"])
		fmt.Println("np: ", float64(e_spin_obs.cos_np[0]) - cosTheta["np"])


		if ievt==0 {
			printEvent(e_partons)
		}
	}
}


// Compute spin-related cosines
func computeSpinCosines(tplus, tminus, lplus, lminus lv.FourVec) (map[string]float64){

	// Get the proper basis
	k, r, n := getSpinBasis(tplus, tminus)

	// Get 3-vectors of lminus (lplus) in tplus (tmius) rest-frame 
	b_to_tplus := tplus.GetBoost()
	lminus_topRF := lminus.ApplyBoost(b_to_tplus).Pvec
	b_to_tminus := tminus.GetBoost()
	lplus_topRF := lplus.ApplyBoost(b_to_tminus).Pvec

	// Fill the six cosines
	getCos := func(a, b r3.Vector, m float64) (float64){ return math.Cos(a.Angle(b.Mul(m)).Radians()) }
	cosTheta := map[string]float64{
		"k+": getCos(lplus_topRF , k,  1),
		"r+": getCos(lplus_topRF , r,  1),
		"n+": getCos(lplus_topRF , n,  1),
		"k-": getCos(lminus_topRF, k, -1),
		"r-": getCos(lminus_topRF, r, -1),
		"n-": getCos(lminus_topRF, n, -1),
	}
	
	return cosTheta
}

// Get spin basis
func getSpinBasis(t, tbar lv.FourVec) (k, r, n r3.Vector) {

	// ttbar rest frame
	ttbar := t.Add(tbar)
	boost_to_ttbar := ttbar.GetBoost()

	// Get top direction in ttbar rest frame
	top_rest := t.ApplyBoost(boost_to_ttbar)
	k = top_rest.Pvec
	k = k.Normalize()
	
	// Get the beam axis (Oz) and coeff to build ortho-normal basis
	beam_axis := r3.Vector{0, 0, 1}
	yval := k.Dot(beam_axis)
	ysign := yval/math.Abs(yval)
	rval := math.Sqrt(1-yval*yval)
	
	// Get n axis: n = sign(y)*r*(beam cross k)  
	n = beam_axis.Cross(k)
	n = n.Mul(rval*ysign)
	
	// Get r axis: r = sign(y)*r*(beam -y*k)
	r = beam_axis.Add(k.Mul(-yval))
	r = r.Mul(1./rval)
	r = r.Mul(ysign)

	return k, r, n
}

// Helper to define the angular-related variables to load
func getAngularVariables(e *SpinObservables) (vars []rtree.ScanVar) {
	vars = []rtree.ScanVar{
		{Name: "kvec"   , Value: &e.kVec},
		{Name: "rvec"   , Value: &e.rVec},
		{Name: "nvec"   , Value: &e.nVec},
		{Name: "dphi_ll", Value: &e.dphi_ll},
		{Name: "dphi_ll", Value: &e.dphi_ll},
		{Name: "dphi_ll", Value: &e.dphi_ll},
		{Name: "cosO_km", Value: &e.cos_km},
		{Name: "cosO_kp", Value: &e.cos_kp},
		{Name: "cosO_rm", Value: &e.cos_rm},
		{Name: "cosO_rp", Value: &e.cos_rp},
		{Name: "cosO_nm", Value: &e.cos_nm},
		{Name: "cosO_np", Value: &e.cos_np},
	}		
	return 
}

// Helper to define the parton-related variables to load
func getPartonVariables(e *PartonEvent) (vars []rtree.ScanVar) {
	vars = []rtree.ScanVar{
		{Name: "t_pt", Value: &e.t.pt},
		{Name: "t_eta", Value: &e.t.eta},
		{Name: "t_phi", Value: &e.t.phi},
		{Name: "t_pid", Value: &e.t.pid},
		{Name: "b_pt", Value: &e.b.pt},
		{Name: "b_eta", Value: &e.b.eta},
		{Name: "b_phi", Value: &e.b.phi},
		{Name: "b_pid", Value: &e.b.pid},
		{Name: "W_pt", Value: &e.W.pt},
		{Name: "W_eta", Value: &e.W.eta},
		{Name: "W_phi", Value: &e.W.phi},
		{Name: "W_pid", Value: &e.W.pid},
		{Name: "l_pt", Value: &e.l.pt},
		{Name: "l_eta", Value: &e.l.eta},
		{Name: "l_phi", Value: &e.l.phi},
		{Name: "l_pid", Value: &e.l.pid},
		{Name: "v_pt", Value: &e.v.pt},
		{Name: "v_eta", Value: &e.v.eta},
		{Name: "v_phi", Value: &e.v.phi},
		{Name: "v_pid", Value: &e.v.pid},
	}		
	return 
}

// Event printing
func printEvent(e PartonEvent) {
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
