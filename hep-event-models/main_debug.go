package main

import (
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

// Event model
type Event struct {
	Jets_Px   [2]float64
	Jets_Py   [2]float64
	Jets_Pz   [2]float64
	Jets_E    [2]float64
	Jets_EMf  [2]float64
	Jets_Ntrk [2]int64
}

func main() {

	var (
		ifname = "TwoEventModels.root"
		tname = "TreeEventArray"
	)
	
	// Open ROOT file
	f, err := groot.Open(ifname)
	if err != nil {
		fmt.Errorf("could not create ROOT file %q: %w", ifname, err)
	}
	defer f.Close()

	// Open the TTree
	obj, err := f.Get(tname)
	if err != nil {
		fmt.Errorf("could not retrieve tree %q: %+v", tname, err)
	}
	tree := obj.(rtree.Tree)
	fmt.Println("Tree is well loaded and has", tree.Entries(), "entries")
	
	// Prepare event reading
	var e Event
	rvars := []rtree.ScanVar{
		{Name: "jets_px"   , Value: &e.Jets_Px  },
		{Name: "jets_py"   , Value: &e.Jets_Py  },
		{Name: "jets_pz"   , Value: &e.Jets_Pz  },
		{Name: "jets_e"    , Value: &e.Jets_E   },
		{Name: "jets_emf"  , Value: &e.Jets_EMf },
		{Name: "jets_ntrks", Value: &e.Jets_Ntrk},
	}
	sc, err := rtree.NewScannerVars(tree, rvars...)
	if err != nil {
		fmt.Errorf("could not create scanner: %+v", err)
	}
	defer sc.Close()
	
	// Event loop
	for sc.Next() {
		
		// Load variable of the event
		ievt, err := sc.Entry(), sc.Scan()
		if err != nil {
			fmt.Errorf("could not scan entry %d: %+v", ievt, err)
		}
		
	}
}
