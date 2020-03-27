// Convert LHE file into a ROOT TTree
package main

import(
	"flag"
	"fmt"
	"io"
	"os"
	
	"go-hep.org/x/hep/lhef"

	_ "go-hep.org/x/hep/groot"
        _ "go-hep.org/x/hep/groot/rtree"
)

func main(){

	// Input arguments
	var (
		ifname = flag.String("f", "ttbar_0j_parton.lhe", "Path to the input LHE file")
	)
	flag.Parse()
	
	// Load LHE file
	f, err := os.Open(*ifname)
	if err != nil {
		panic(err)
	}

	// Get LHE decoder
	lhedec, err := lhef.NewDecoder(f)
	if err != nil {
		panic(err)
	}

	// Loop over events
	for i := 0; ; i++ {

		// Decode this event, stop if the end of file is reached
		evt, err := lhedec.Decode()
		if err == io.EOF {
			break
		}

		// Fill the event
		fmt.Println(*evt)
		
		// Write the TTree
	}

	
	
}