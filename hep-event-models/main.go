// Manipulate two different structures of input events using go interface 
package main

import (
	"math"
	"flag"
	"strings"
	"log"
	
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {

	var (
		fname  = flag.String("f"  , "TwoEventModels.root", "Input file name")
		emodel = flag.String("edm", "flat", "Input event model: 'flat', 'array'")
	)

	flag.Parse()
	
	eventLoop(*fname, *emodel)
	eventLoop(*fname, *emodel)
	
}

// Perform the event loop
func eventLoop(ifname, emodel string) {
	
	// Open ROOT file
	f, err := groot.Open(ifname)
	if err != nil {
		log.Fatalf("could not create ROOT file %q: %w", ifname, err)
	}
	defer f.Close()

	// Choose the (tree name, event model) depending on the specified event model
	var tname string
	var eIn EventIn
        switch strings.ToLower(emodel) {
	case "flat":
		tname = "TreeEventFlat"
		eIn = &EventInFlat{}
	case "array":
		tname = "TreeEventArray"
		eIn = &EventInArray{}
	}
	
	// Open the TTree
	obj, err := f.Get(tname)
	if err != nil {
		log.Fatalf("could not retrieve tree %q: %+v", tname, err)
	}
	tree := obj.(rtree.Tree)
	
	// Prepare event reading
	rvars := eIn.GetTreeScannerVars()
	sc, err := rtree.NewScannerVars(tree, rvars...)
	if err != nil {
		log.Fatalf("could not create scanner: %+v", err)
	}
	defer sc.Close()
	
	// Prepare the output file, tree and event processing
	fnameOut := "ProcessedFromEvent" + strings.Title(emodel) + ".root"
	fout, err := groot.Create(fnameOut)
	if err != nil {
		log.Fatalf("could not create ROOT file %q: %w", fnameOut, err)
	}
	defer fout.Close()
	var eOut EventOut
	wvars := eOut.GetTreeWriterVars()
	treeOut, err := rtree.NewWriter(fout, "TreeEventOut", wvars)
	if err != nil {
		log.Fatalf("could not create scanner: %+v", err)
	}
	defer treeOut.Close()

	// Event loop
	for sc.Next() {
		
		// Load variable of the event
		ievt, err := sc.Entry(), sc.Scan()
		if err != nil {
			log.Fatalf("could not scan entry %d: %+v", ievt, err)
		}
		
		// Copy the 'input' event format into the 'output' event format
		eIn.CopyTo(&eOut)
		
		// Process the event
		eOut.Process()

		// Write output tree
		_, err = treeOut.Write()
		if err != nil {
			log.Fatalf("could not write output tree %d: %+v", ievt, err)
		}
	}
}

// Output event model, structured in two RecoJet objects,
// with additionnal information.
type EventOut struct {
	Jet1, Jet2 RecoJet
	InvMass float64
}

// RecoJet definition grouping the input information
// and additional variables in the same structure
type RecoJet struct {
	Px, Py, Pz, E float64
	Ntrk int64
	EMf float64
	HADf float64
}

// Function returning a writer to associate each variable to branch name
// In this example, only the information of interested in saved.
func (e *EventOut) GetTreeWriterVars() []rtree.WriteVar {
	return []rtree.WriteVar{
		{Name: "jet_mass" , Value: &e.InvMass  },
		{Name: "jet1_hadf", Value: &e.Jet1.HADf},
		{Name: "jet2_hadf", Value: &e.Jet2.HADf},		
	}
}

// Function implementing the processing of EventOut
func (e *EventOut) Process() EventOut {

	// Compute M(j,j)
	E2 := e.Jet1.E*e.Jet1.E
	P2 := e.Jet1.Px*e.Jet1.Px+ e.Jet1.Py*e.Jet1.Py + e.Jet1.Pz*e.Jet1.Pz
	e.InvMass = math.Sqrt(math.Abs(E2 - P2))

	// Compute the hadronic fraction of each jet
	e.Jet1.HADf = 1 - e.Jet1.EMf
	e.Jet2.HADf = 1 - e.Jet2.EMf

	return *e
}

// Interface for generic input event model which must be
// 1. read, ie associate a name to a value: GetTreeScannerVars()
// 2. converted into EventOut: CopyTo() 
type EventIn interface {
	GetTreeScannerVars() []rtree.ScanVar
	CopyTo(evt *EventOut)
}

// Input Event model made of 12 (flat) numbers
type EventInFlat struct {
	Jet1_Px  , Jet2_Px   float64
	Jet1_Py  , Jet2_Py   float64
	Jet1_Pz  , Jet2_Pz   float64
	Jet1_E   , Jet2_E    float64
	Jet1_EMf , Jet2_EMf  float64
	Jet1_Ntrk, Jet2_Ntrk int64
}

// Implementation of the reading of EventInFlat
func (eIn *EventInFlat) GetTreeScannerVars() []rtree.ScanVar {
	return []rtree.ScanVar{

		// Jet 1
		{Name: "jet_px_1"  , Value: &eIn.Jet1_Px  },
		{Name: "jet_py_1"  , Value: &eIn.Jet1_Py  },
		{Name: "jet_pz_1"  , Value: &eIn.Jet1_Pz  },
		{Name: "jet_e_1"   , Value: &eIn.Jet1_E   },
		{Name: "jet_ntrk_1", Value: &eIn.Jet1_Ntrk},
		{Name: "jet_emf_1" , Value: &eIn.Jet1_EMf },
		
		// Jet 2
		{Name: "jet_px_2"  , Value: &eIn.Jet2_Px  },
		{Name: "jet_py_2"  , Value: &eIn.Jet2_Py  },
		{Name: "jet_pz_2"  , Value: &eIn.Jet2_Pz  },
		{Name: "jet_e_2"   , Value: &eIn.Jet2_E   },
		{Name: "jet_ntrk_2", Value: &eIn.Jet2_Ntrk},
		{Name: "jet_emf_2" , Value: &eIn.Jet2_EMf },
	}
}

// Implementation of copying EventInFlat to EventOut
func (eIn EventInFlat) CopyTo(eOut *EventOut) {

	eOut.Jet1 = RecoJet{
		Px: eIn.Jet1_Px,
		Py: eIn.Jet1_Py,
		Pz: eIn.Jet1_Pz,
		E: eIn.Jet1_E,
		Ntrk: eIn.Jet1_Ntrk,
		EMf: eIn.Jet1_EMf,
	}

	eOut.Jet2 = RecoJet{
		Px: eIn.Jet2_Px,
		Py: eIn.Jet2_Py,
		Pz: eIn.Jet2_Pz,
		E: eIn.Jet2_E,
		Ntrk: eIn.Jet2_Ntrk,
		EMf: eIn.Jet2_EMf,
	}
}


// Input event model made of 6 arrays of 2-elements
type EventInArray struct {
	Jets_Px   [2]float64
	Jets_Py   [2]float64
	Jets_Pz   [2]float64
	Jets_E    [2]float64
	Jets_EMf  [2]float64
	Jets_Ntrk [2]int64
}

// Implementation of the reading of EventInArray
func (eIn *EventInArray) GetTreeScannerVars() []rtree.ScanVar {
	return []rtree.ScanVar{
		{Name: "jets_px"  , Value: &eIn.Jets_Px  },
		{Name: "jets_py"  , Value: &eIn.Jets_Py  },
		{Name: "jets_pz"  , Value: &eIn.Jets_Pz  },
		{Name: "jets_e"   , Value: &eIn.Jets_E   },
		{Name: "jets_ntrk", Value: &eIn.Jets_Ntrk},
		{Name: "jets_emf" , Value: &eIn.Jets_EMf },
	}
}

// Implementation of copying EventInFlat to EventOut
func (eIn EventInArray) CopyTo(eOut *EventOut) {
	var reco_jets [2]RecoJet
	for i := range eIn.Jets_Px {
		reco_jets[i] = RecoJet{
			Px: eIn.Jets_Px[i],
			Py: eIn.Jets_Py[i],
			Pz: eIn.Jets_Pz[i],
			E: eIn.Jets_E[i],
			Ntrk: eIn.Jets_Ntrk[i],
			EMf: eIn.Jets_EMf[i],
		}
	}
	eOut.Jet1 = reco_jets[0]
	eOut.Jet2 = reco_jets[1]
}

