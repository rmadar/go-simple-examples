package main

import (
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
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

type Event struct {
	t    Particles
	b    Particles
	W    Particles
	l    Particles
	v    Particles
}

type Particles struct {
	pt  []float32
	eta []float32
	phi []float32
	m   []float32
	pid []int32
}

// Event loop
func eventLoop(fname string, tname string, evtmax int64) {

	// Open the root file and get the tree
	fmt.Println("Opening TTree", tname, "in ROOT file", fname)
	file := openRootFile(fname)
	tree := getTtree(file, tname)

	// Event and variables to load
	var e Event
	vars := []rtree.ScanVar{
		{Name: "t_pt" , Value: &e.t.pt},
		{Name: "t_eta", Value: &e.t.eta},
		{Name: "t_phi", Value: &e.t.phi},
		{Name: "t_pid", Value: &e.t.pid},
	}

	// Create a scanner to perform the event loop
	sc, _ := rtree.NewScannerVars(tree, vars...)
	
	// Actual event loop
	for sc.Next() && sc.Entry() < evtmax {
		iev := sc.Entry()

		if iev%1000==0 {
			fmt.Println("Evt:", iev)
			printEvent(e)
		}
	}
}

// Event printing
func printEvent(e Event){
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
}

// Helper to open a ROOT file
func openRootFile(fname string) (*riofs.File) {
	f, _ := groot.Open(fname)
	return f
}

// Helper to get a TTree
func getTtree(f *riofs.File, tname string) (rtree.Tree){
	obj, _ := f.Get(tname)
	return obj.(rtree.Tree)
}

