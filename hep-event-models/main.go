
// Exemple of how to use go interface to define several HEP
// event models. The goal is to have a single event model
// at the user level, and handle the differences & conversion
// in one package using interfaces. This example is based on
// ROOT object translated in go in groot (go-hep package),
// but can be easily generalized.
//
// The use case in mind for this example is to have two inputs
// formats (say 'EventInFlat' and 'EventInArrays') for one unique
// output format (say 'EventOut'). The user would have to define
// once. The format of 'EventInFlat' and 'EventInFlat', and how to
// translate each of them into 'EventOut'. Then, with a simple flag
// (passed as argument) one can run the same executable on the two
// input format. In summary, we have:
//  1. EventOut, to be processed (add new variables) and saved
//  2. EventInFlat, to be loaded and converted into EventOut
//  3. EventInArray, to be loaded and converted into EventOut
//
// For this example, we consider events with 2-jets in the final
// state, described by 12 numbers, 2 x (4-vectors + number of
// tracks + EM fraction), that can be organized in 3 different ways.

package main

import (
	"fmt"
	"math"
	
	"go-hep.org/x/hep/groot/rtree"
)




func main() {
	
	fmt.Println("Testing 'go interface' concept")
	
	
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

// Input event model made of 6 arrays of 2-elements
type EventInArray struct {
	Jets_Px   [2]float64
	Jets_Py   [2]float64
	Jets_Pz   [2]float64
	Jets_E    [2]float64
	Jets_EMf  [2]float64
	Jets_Ntrk [2]int64
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
func (eOut *EventOut) GetTreeWriter() []rtree.WriteVar {
	return []rtree.WriteVar{
		{Name: "jet_mass" , Value: &eOut.InvMass  },
		{Name: "jet1_hadf", Value: &eOut.Jet1.HADf},
		{Name: "jet2_hadf", Value: &eOut.Jet2.HADf},		
	}
}

// Function implementing the processing of EventOut
func (e *EventOut) Process() EventOut {

	// Compute M(j,j)
	e.InvMass = math.Sqrt(math.Abs(e.Jet1.E*e.Jet1.E - (
		e.Jet1.Px*e.Jet1.Px
		+ e.Jet1.Py*e.Jet1.Py
		+ e.Jet1.Pz*e.Jet1.Pz)))

	// Compute the hadronic fraction of each jet
	e.Jet1.HADf = 1 - e.Jet1.EMf
	e.Jet2.HADf = 1 - e.Jet2.EMf

	return *e
}

// Interface for generic input event model which must be
// 1. read, ie associate a name to a value: GetTreeScanner()
// 2. converted into EventOut: CopyTo() 
type EventInput interface {
	GetTreeScanner() []rtree.ScanVar
	CopyTo(evt *EventOut)
}

// Implementation of the reading of EventInFlat
func (eIn *EventInFlat) GetTreeScanner() []rtree.ScanVar {
	return []rtree.ScanVar{

		// Jet 1
		{Name: "jet_px_1"  , Value: &eIn.Jet1_Px  },
		{Name: "jet_py_1"  , Value: &eIn.Jet1_Py  },
		{Name: "jet_pz_1"  , Value: &eIn.Jet1_Pz  },
		{Name: "jet_e_1"   , Value: &eIn.Jet1_E   },
		{Name: "jet_ntrk_1", Value: &eIn.Jet1_Ntrk},
		{Name: "jet_emf_1" , Value: &eIn.Jet1_EMf },
		
		// Jet 2
		{Name: "jet_px_2"  , Value: &eIn.Jet1_Px  },
		{Name: "jet_py_2"  , Value: &eIn.Jet1_Py  },
		{Name: "jet_pz_2"  , Value: &eIn.Jet1_Pz  },
		{Name: "jet_e_2"   , Value: &eIn.Jet1_E   },
		{Name: "jet_ntrk_2", Value: &eIn.Jet1_Ntrk},
		{Name: "jet_emf_2" , Value: &eIn.Jet1_EMf },
	}
}

// Implementation of copying EventInFlat to EventOut
func (eIn *EventInFlat) CopyTo(eOut *EventOut) {

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


// Implementation of the reading of EventInArray
func (eIn *EventInArray) GetTreeScanner() []rtree.ScanVar {
	return []rtree.ScanVar{
		{Name: "jets_px"   , Value: &eIn.Jets_Px  },
		{Name: "jets_py"   , Value: &eIn.Jets_Py  },
		{Name: "jets_pz"   , Value: &eIn.Jets_Pz  },
		{Name: "jets_e"    , Value: &eIn.Jets_E   },
		{Name: "jets_ntrks", Value: &eIn.Jets_Ntrk},
		{Name: "jets_emf"  , Value: &eIn.Jets_EMf },
	}
}

// Implementation of copying EventInFlat to EventOut
func (eIn *EventInArray) CopyTo(eOut *EventOut) {
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

