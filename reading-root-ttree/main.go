package main

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {

	var (
		fname  = "ttbar_0j_parton.root"
		tname  = "test"
		evtmax = int64(10000)
	)

	eventLoop(fname, tname, evtmax)
}

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
	dphi_ll []float32
	cos_km []float32
	cos_kp []float32
	cos_rm []float32
	cos_rp []float32
	cos_nm []float32
	cos_np []float32
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
		iev := sc.Entry()

		// Load variable of the event
		err := sc.Scan()
		if err != nil {
			log.Fatalf("could not scan entry %d: %+v", iev, err)
		}

		// Print
		if iev%10000 == 0 {
			fmt.Println("Evt:", iev)
			printEvent(e_partons)
		}
	}
}

// Helper to define the angular-related variables to load
func getAngularVariables(e *SpinObservables) (vars []rtree.ScanVar) {
	vars = []rtree.ScanVar{
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
