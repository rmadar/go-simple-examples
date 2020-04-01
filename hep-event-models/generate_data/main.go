
// Generate ROOT file with two different event models
//  1. flat data
//  2. array data

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)


func main() {

	var (
		fname  = flag.String("f", "../TwoEventModels.root", "path to ROOT file to create")
		evtmax = flag.Int64("n", 10000, "number of events to generate")
	)

	flag.Parse()
	
	generateData(*fname, *evtmax)
}


// Input Event model made of 12 (flat) numbers
type EventFlat struct {
	Jet1_Px  , Jet2_Px   float64
	Jet1_Py  , Jet2_Py   float64
	Jet1_Pz  , Jet2_Pz   float64
	Jet1_E   , Jet2_E    float64
	Jet1_EMf , Jet2_EMf  float64
	Jet1_Ntrk, Jet2_Ntrk int64
}

// Input event model made of 6 arrays of 2-elements
type EventArray struct {
	Jets_Px   [2]float64
	Jets_Py   [2]float64
	Jets_Pz   [2]float64
	Jets_E    [2]float64
	Jets_EMf  [2]float64
	Jets_Ntrk [2]int64
}

// Creating the two TTrees in the same file
func generateData(fname string, evtmax int64) {
	f, err := groot.Create(fname)
	if err != nil {
		fmt.Errorf("could not create ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	// Flat Event format
	var eFlat EventFlat
	wvars_flat := []rtree.WriteVar{

		// Jet 1
		{Name: "jet_px_1"  , Value: &eFlat.Jet1_Px  },
		{Name: "jet_py_1"  , Value: &eFlat.Jet1_Py  },
		{Name: "jet_pz_1"  , Value: &eFlat.Jet1_Pz  },
		{Name: "jet_e_1"   , Value: &eFlat.Jet1_E   },
		{Name: "jet_ntrk_1", Value: &eFlat.Jet1_Ntrk},
		{Name: "jet_emf_1" , Value: &eFlat.Jet1_EMf },
		
		// Jet 2
		{Name: "jet_px_2"  , Value: &eFlat.Jet2_Px  },
		{Name: "jet_py_2"  , Value: &eFlat.Jet2_Py  },
		{Name: "jet_pz_2"  , Value: &eFlat.Jet2_Pz  },
		{Name: "jet_e_2"   , Value: &eFlat.Jet2_E   },
		{Name: "jet_ntrk_2", Value: &eFlat.Jet2_Ntrk},
		{Name: "jet_emf_2" , Value: &eFlat.Jet2_EMf },
	}

	treeFlat, err := rtree.NewWriter(f, "TreeEventFlat", wvars_flat)
	if err != nil {
		fmt.Errorf("could not create tree writer: %w", err)
	}
	defer treeFlat.Close()
	
	// Array Event format
	var eArray EventArray
	wvars_array := []rtree.WriteVar{
		{Name: "jets_px"  , Value: &eArray.Jets_Px  },
		{Name: "jets_py"  , Value: &eArray.Jets_Py  },
		{Name: "jets_pz"  , Value: &eArray.Jets_Pz  },
		{Name: "jets_e"   , Value: &eArray.Jets_E   },
		{Name: "jets_ntrk", Value: &eArray.Jets_Ntrk},
		{Name: "jets_emf" , Value: &eArray.Jets_EMf },
	}

	treeArray, err := rtree.NewWriter(f, "TreeEventArray", wvars_array)
	if err != nil {
		fmt.Errorf("could not create tree writer: %w", err)
	}
	defer treeArray.Close()

	var P2 float64
	for i := int64(0); i < evtmax; i++ {

		// Event flat - jet 1
		eFlat.Jet1_Px = rand.Float64()
		eFlat.Jet1_Py = rand.Float64()
		eFlat.Jet1_Pz = rand.Float64()
		P2 = eFlat.Jet1_Px*eFlat.Jet1_Px + eFlat.Jet1_Py*eFlat.Jet1_Py + eFlat.Jet1_Pz*eFlat.Jet1_Pz
		eFlat.Jet1_E  = math.Sqrt(P2) + rand.Float64()*5
		eFlat.Jet1_Ntrk = rand.Int63n(10)
		eFlat.Jet1_EMf = rand.Float64()
		// Event flat - jet 2
		eFlat.Jet2_Px = rand.Float64()
		eFlat.Jet2_Py = rand.Float64()
		eFlat.Jet2_Pz = rand.Float64()
		P2 = eFlat.Jet2_Px*eFlat.Jet2_Px + eFlat.Jet2_Py*eFlat.Jet2_Py + eFlat.Jet2_Pz*eFlat.Jet2_Pz
		eFlat.Jet2_E  = math.Sqrt(P2) + rand.Float64()*2
		eFlat.Jet2_Ntrk = rand.Int63n(10)
		eFlat.Jet2_EMf = rand.Float64()
		// Event flat - write the etree
		_, err = treeFlat.Write()
		if err != nil {
			fmt.Errorf("could not write event %d: %w", i, err)
		}

		// Event Array - jet 1
		eArray.Jets_Px[0] = eFlat.Jet1_Px 
		eArray.Jets_Py[0] = eFlat.Jet1_Py 
		eArray.Jets_Pz[0] = eFlat.Jet1_Pz 
		eArray.Jets_E[0] = eFlat.Jet1_E
		eArray.Jets_Ntrk[0] = eFlat.Jet1_Ntrk 
		eArray.Jets_EMf[0] = eFlat.Jet1_EMf
		// Event Array - jet 2
		eArray.Jets_Px[1] = eFlat.Jet2_Px 
		eArray.Jets_Py[1] = eFlat.Jet2_Py 
		eArray.Jets_Pz[1] = eFlat.Jet2_Pz 
		eArray.Jets_E[1] = eFlat.Jet2_E
		eArray.Jets_Ntrk[1] = eFlat.Jet2_Ntrk 
		eArray.Jets_EMf[1] = eFlat.Jet2_EMf 
		// Event array - write the tree
		_, err = treeArray.Write()
		if err != nil {
			fmt.Errorf("could not write event %d: %w", i, err)
		}
		
	}

	err = treeFlat.Close()
	if err != nil {
		fmt.Errorf("could not close tree writer: %w", err)
	}

	err = treeArray.Close()
	if err != nil {
		fmt.Errorf("could not close tree writer: %w", err)
	}
	
	err = f.Close()
	if err != nil {
		fmt.Errorf("could not close ROOT file: %w", err)
	}

}
