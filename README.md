# Simple Go examples

This repository contains few examples in go, mostly connected to high energy physics (HEP). Mostly:
  + `CLs`: computation of a stat-only CLs exclusion [wikipedia](https://en.wikipedia.org/wiki/CLs_method_(particle_physics))	
  + `plotting`: very short example of plots produced in go
  + `lhe2root`: example of how to convert a LHE file into a ROOT file
  + `reading-root-ttree`: event loop based on a `ROOT::TTree`, the main HEP software.


### CLs exclusion

### LHE to ROOT - based on [go-hep](https://go-hep.org/)

LHE format is convention to store data from particle collision into an ASCI file. A LHE parser is available in [go-hep](https://godoc.org/go-hep.org/x/hep/lhef) and is used to create a `TTree` for a 10000 proton-proton collisions leading to a top-antitop quark pair production.

### Reading a `TTree` - based on [go-hep](https://go-hep.org/)

In this example, the initial `TTree` - stored in [ttbar_0j_parton.root](reading-root-ttree/main.go) - was produced from a LHE file [[arXiv:0609.017](https://arxiv.org/abs/hep-ph/0609017)] describing 10000 proton-proton collisions leading to a top-antitop quark pair production, as predicted by MadGraph tool [[arXiv:1405.0301](https://arxiv.org/abs/1405.0301)], ran at the leading order.
These collisions are described at the parton level only and each event is described by
  + partonic intial state: parton flavour and momentum
  + partonic final state: 4-vectors for each particle in the decay `t->Wb->lvb`

The program [reading-root-ttree/main.go](reading-root-ttree/main.go) loads some variables of the `TTree`, compute
some angular variables probing the spin correlation between the top and the antitop quarks [e.g. [arXiv:1612.07004](https://arxiv.org/abs/1612.07004)]. These involves Lorentz transformation and simple geometrical calculations, and this progam relies then on the [lorentzvector package](https://godoc.org/github.com/rmadar/go-lorentz-vector/lv). The commands
```bash
cd reading-root-ttree
go run ./main.go
```
will produce a ROOT file containing the new `TTree` with 10 variables:
  + `dphi_ll`: lab-frame angle between the two leptons
  + `k, r, n`: the three 3-vectors of the spin basis
  + `cos(Theta[axis, lepton])`: 6 cosines for 3 axis and 2 leptons