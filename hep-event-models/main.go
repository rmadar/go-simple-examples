
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
//  1. EventOut, to be saved
//  2. EventInFlat, to be loaded and converted into EventOut
//  3. EventInArray, to be loaded and converted into EventOut
//
// For this example, we consider events with 2-jets in the final
// state, described by 12 numbers, 2 x (4-vectors + number of
// tracks + EM fraction), that can be organized in 3 different ways.

package main

import (
	"fmt"

	"go-hep.org/x/hep/groot/rtree"
)

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
	Jet_px   [2]float64
	Jet_py   [2]float64
	Jet_pz   [2]float64
	Jet_e    [2]float64
	Jet_emf  [2]float64
	Jet_ntrk [2]int64
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
	Ntrks int64
	EMf float64
	HADf float64
}

// Function returning a writer to associate each variable to branch name
// In this example, only the information of interested in saved.
func (*e EventOut) GetTreeWriter() []rtree.WriteVar {
	return []rtree.WriteVar{
		{Name: "jet_mass" , Value: &e.InvMass},
		{Name: "jet1_hadf", Value: &e.Jet1.HADf},
		{Name: "jet2_hadf", Value: &e.Jet2.HADf},		
	}
}

// Interface for generic input event model which must be
// 1. read, ie associate a name to a value: GetTreeScanner()
// 2. converted into EventOut: CopyTo() 
type EventInput interface {
	GetTreeScanner() []rtree.ScanVar
	CopyTo(evt *EventOut)
}

// Implementation of reading part for EventInFlat
func (e *EventInFlat) GetTreeScanner() []rtree.ScanVar {
	
}


func main() {
	
	fmt.Println("Testing 'go interface' concept")
}
