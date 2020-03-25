// Exemple of reading a ROOT TTree computing some spin-correlation observables
package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"github.com/golang/geo/r3"
	"github.com/rmadar/go-lorentz-vector/lv"
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
		fname  = flag.String("f", "ttbar_0j_parton.root", "path to ROOT file to analyze")
		tname  = flag.String("t", "spinCorrelation", "ROOT Tree name to analyze")
		evtmax = flag.Int64("n", 10000, "number of events to analyze")
	)

	flag.Parse()

	eventLoop(*fname, *tname, *evtmax)
}

// Event loop
func eventLoop(fname string, tname string, evtmax int64) {

	// Open the root file and get the tree
	fmt.Println("Opening TTree", tname, "in ROOT file", fname)
	file := openRootFile(fname)
	tree := getTtree(file, tname)
	
	// Create a scanner to perform the event loop on the input tree
	var (
		e_partons PartonEvent
		e_spin_obs SpinObservables
		rvars = make([]rtree.ScanVar, 0)
	)
	rvars = append(rvars, getReadPartonVariables(&e_partons)...)
	rvars = append(rvars, getReadAngularVariables(&e_spin_obs)...)
	sc, err := rtree.NewScannerVars(tree, rvars...)
	if err != nil {
		log.Fatalf("could not create scanner: %+v", err)
	}
	defer sc.Close()
	
	// Ceate a new file, new writer to save new variables in a tree
	fnameOut := "test.root"
	fout, err := groot.Create(fnameOut)
	if err != nil {
		log.Fatalf("could not create ROOT file %q: %w", fnameOut, err)
	}
	defer fout.Close()
	var (
		e_spin_obsOut SpinObservables
		//e_spin_obsCopy SpinObservables
		nUnit_size int32
		nThree_size int32
	)
	// this structure of slices is not very good but is necessary
	// to keep the same 'SpinObservables' structure for reading/writting.
	// (this is due to how the input tree is produced)
	wvars := []rtree.WriteVar{ 

		// Counter for slices
		{Name: "N1"     , Value: &nUnit_size},
		{Name: "N3"     , Value: &nThree_size},

		// Locally computed variables
		{Name: "kvec"   , Value: &e_spin_obsOut.kVec, Count: "N3"},
		{Name: "rvec"   , Value: &e_spin_obsOut.rVec, Count: "N3"},
		{Name: "nvec"   , Value: &e_spin_obsOut.nVec, Count: "N3"},
		{Name: "dphi_ll", Value: &e_spin_obsOut.dphi_ll, Count: "N1"},
		{Name: "cosO_km", Value: &e_spin_obsOut.cos_km, Count: "N1"},
		{Name: "cosO_kp", Value: &e_spin_obsOut.cos_kp, Count: "N1"},
		{Name: "cosO_rm", Value: &e_spin_obsOut.cos_rm, Count: "N1"},
		{Name: "cosO_rp", Value: &e_spin_obsOut.cos_rp, Count: "N1"},
		{Name: "cosO_nm", Value: &e_spin_obsOut.cos_nm, Count: "N1"},
		{Name: "cosO_np", Value: &e_spin_obsOut.cos_np, Count: "N1"},

		// Copy of original variables
		{Name: "kvec_orig"   , Value: &e_spin_obs.kVec, Count: "N3"},
		{Name: "rvec_orig"   , Value: &e_spin_obs.rVec, Count: "N3"},
		{Name: "nvec_orig"   , Value: &e_spin_obs.nVec, Count: "N3"},
		{Name: "dphi_ll_orig", Value: &e_spin_obs.dphi_ll, Count: "N1"},
		{Name: "cosO_km_orig", Value: &e_spin_obs.cos_km, Count: "N1"},
		{Name: "cosO_kp_orig", Value: &e_spin_obs.cos_kp, Count: "N1"},
		{Name: "cosO_rm_orig", Value: &e_spin_obs.cos_rm, Count: "N1"},
		{Name: "cosO_rp_orig", Value: &e_spin_obs.cos_rp, Count: "N1"},
		{Name: "cosO_nm_orig", Value: &e_spin_obs.cos_nm, Count: "N1"},
		{Name: "cosO_np_orig", Value: &e_spin_obs.cos_np, Count: "N1"},
	}
	tout, err := rtree.NewWriter(fout, "newSpinObs", wvars)
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
		if ievt==0 {
			printEvent(e_partons)
		}
		
		// Getting slice of particles
		tops := e_partons.t
		leptons := e_partons.l
				
		// Re-computing spin observables
		var (
			loadVec = lv.NewFourVecPtEtaPhiM
			tplus_P4 = loadVec(float64(tops.pt[0]), float64(tops.eta[0]), float64(tops.phi[0]), float64(tops.m[0]))
			tminus_P4 = loadVec(float64(tops.pt[1]), float64(tops.eta[1]), float64(tops.phi[1]), float64(tops.m[1]))
			lplus_P4 = loadVec(float64(leptons.pt[0]), float64(leptons.eta[0]), float64(leptons.phi[0]), float64(leptons.m[0]))
			lminus_P4 = loadVec(float64(leptons.pt[1]), float64(leptons.eta[1]), float64(leptons.phi[1]), float64(leptons.m[1]))
		)
		cosTheta := computeSpinCosines(tplus_P4, tminus_P4, lplus_P4, lminus_P4)

		// Compare with stored variables
		fmt.Println("\nComparing stored and recomputed spin variables:")
		fmt.Println("k-: ", e_spin_obs.cos_km[0], cosTheta["k-"])
		fmt.Println("r-: ", e_spin_obs.cos_rm[0], cosTheta["r-"])
		fmt.Println("n-: ", e_spin_obs.cos_nm[0], cosTheta["n-"])
		fmt.Println("k+: ", e_spin_obs.cos_kp[0], cosTheta["k+"])
		fmt.Println("r+: ", e_spin_obs.cos_rp[0], cosTheta["r+"])
		fmt.Println("n+: ", e_spin_obs.cos_np[0], cosTheta["n+"])

		// Compare spin basis vectors
		k, r, n := getSpinBasis(tplus_P4, tminus_P4)
		getVector := func (x[]float32) (r3.Vector) {return r3.Vector{float64(x[0]), float64(x[1]), float64(x[2])} }
		k_ref, r_ref, n_ref := getVector(e_spin_obs.kVec), getVector(e_spin_obs.rVec), getVector(e_spin_obs.nVec)
		fmt.Println(k.Add(k_ref.Mul(-1)))
		fmt.Println(r.Add(r_ref.Mul(-1)))
		fmt.Println(n.Add(n_ref.Mul(-1)))

		// Compare four-vectors obtained with fmom
		lplus_fmom := fmom.NewPtEtaPhiM(float64(leptons.pt[0]), float64(leptons.eta[0]), float64(leptons.phi[0]), float64(leptons.m[0]))
		fmt.Printf("(px, py, pz, E)_fmom = (%.2f, %.2f, %.2f, %.2f)\n", lplus_fmom.Px(), lplus_fmom.Py(), lplus_fmom.Pz(), lplus_fmom.E() )
		fmt.Printf("(px, py, pz, E)_lv   = (%.2f, %.2f, %.2f, %.2f)\n", lplus_P4.Px(), lplus_P4.Py(), lplus_P4.Pz(), lplus_P4.E() )
		fmt.Printf("Phi[fmom, lv] = [%.2f, %.2f]\n", lplus_fmom.Phi(), lplus_P4.Phi())

		// Save the newly computed info into a TTree
		r3Vec2Slice := func(v r3.Vector) ([]float32) {return []float32 {float32(v.X), float32(v.Y), float32(v.Z)}}
		nUnit_size, nThree_size = 1, 3
		e_spin_obsOut.kVec = r3Vec2Slice(k)
		e_spin_obsOut.rVec = r3Vec2Slice(r)
		e_spin_obsOut.nVec = r3Vec2Slice(n)
		e_spin_obsOut.dphi_ll = []float32 {float32(lplus_P4.DeltaPhi(lminus_P4))}
		e_spin_obsOut.cos_km  = []float32 {float32(cosTheta["k-"])}
		e_spin_obsOut.cos_rm  = []float32 {float32(cosTheta["r-"])}
		e_spin_obsOut.cos_nm  = []float32 {float32(cosTheta["n-"])}
		e_spin_obsOut.cos_kp  = []float32 {float32(cosTheta["k+"])}
		e_spin_obsOut.cos_rp  = []float32 {float32(cosTheta["r+"])}
		e_spin_obsOut.cos_np  = []float32 {float32(cosTheta["n+"])}
		tout.Write()
	}
}


// Compute spin-related cosines
func computeSpinCosines(tplus, tminus, lplus, lminus lv.FourVec) (map[string]float64){

	// Get the proper basis
	k, r, n := getSpinBasis(tplus, tminus)

	// Get 3-vectors of lminus (lplus) in tplus (tmius) rest-frame 
	lplusRF := lplus.ToRestFrameOf(tplus).Pvec
	lminusRF := lminus.ToRestFrameOf(tminus).Pvec

	// Fill the six cosines
	getCos := func(a, b r3.Vector, m float64) (float64){
		// an, bn := a.Norm(), b.Norm()
		// return a.Dot(b) / (an*bn)
		return math.Cos(a.Angle(b.Mul(m)).Radians())
		//return 2.0
	}
	cosTheta := map[string]float64{
		"k+": getCos(lplusRF , k,  1),
		"r+": getCos(lplusRF , r,  1),
		"n+": getCos(lplusRF , n,  1),
		"k-": getCos(lminusRF, k, -1),
		"r-": getCos(lminusRF, r, -1),
		"n-": getCos(lminusRF, n, -1),
	}
	
	return cosTheta
}

// Get spin basis
func getSpinBasis(t, tbar lv.FourVec) (k, r, n r3.Vector) {

	// ttbar rest frame
	ttbar := t.Add(tbar)
	boost_to_ttbar := ttbar.GetBoost()

	// Get top direction in ttbar rest frame
	top_rest := t.ApplyBoost(boost_to_ttbar.Mul(-1))
	k = top_rest.Pvec
	k = k.Normalize()
	
	// Get the beam axis (Oz) and coeff to build ortho-normal basis
	beam_axis := r3.Vector{0, 0, 1}
	yval := k.Dot(beam_axis)
	ysign := yval/math.Abs(yval)
	rval := math.Sqrt(1-yval*yval)

	// Get r axis: r = sign(y)/r * (beam -y*k)
	r = beam_axis.Add(k.Mul(-yval))
	r = r.Mul(1./rval*ysign)
	
	// Get n axis: n = sign(y)/r * (beam cross k)  
	n = beam_axis.Cross(k)
	n = n.Mul(1./rval*ysign)
	
	return k, r, n
}

// Helper to define the angular-related variables to load
func getReadAngularVariables(e *SpinObservables) (vars []rtree.ScanVar) {
	vars = []rtree.ScanVar{
		{Name: "kvec"   , Value: &e.kVec},
		{Name: "rvec"   , Value: &e.rVec},
		{Name: "nvec"   , Value: &e.nVec},
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
func getReadPartonVariables(e *PartonEvent) (vars []rtree.ScanVar) {
	vars = []rtree.ScanVar{
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
