# Different event-models at once

Exemple of how to use go interface to define several HEP event models. The goal is to have a single event model at the user level, and handle the differences & conversion in thanks to an interface. This example is based on ROOT object translated in go *via* [go-hep](https://go-hep.org/), but can be easily generalized.

The use case in mind for this example is to have two inputs event formats, `EventInFlat` and `EventInArrays`, for one unique output format `EventOut`. The user just has to define **once**, in an interface, the format of `EventInFlat` and `EventInFlat`, and how to translate each of them into the format `EventOut`. Then, with a simple flag (passed as argument) one can run the same executable on different input event model in *a perfectly transparent way*

In summary, we have the following workflow:
```
 ----------------------       ----------       ------------------       ----------------
| *Interface EventIn*  |     |          |     |                  |     |                |
| EventInFlat reading  | --> | EventOut | --> | Event processing | --> | Output writing |
| EventInArray reading |     |          |     |                  |     |                |
 ----------------------       ----------       ------------------       ----------------
```
where `EventInFlat` and `EventInArray` statisfy the same **`EventIn` interface** defined:
```go
type EventIn interface {
     GetTreeScannerVars() []rtree.ScanVar // TTree reading (contains branch name <-> variable association)
     CopyTo(evt *EventOut)                // Convertion to EventOut event model
}			     
```

For this example, we consider events with 2-jets in the final state, described by 12 numbers, 2 x (4-vectors + number of tracks + EM fraction), that can be organized in 2 different ways. The `EventOut` onlyt store information of interest, and add new observables.

```go
// Flat tree input 
type EventInFlat struct {
	Jet1_Px  , Jet2_Px   float64
	Jet1_Py  , Jet2_Py   float64
	Jet1_Pz  , Jet2_Pz   float64
	Jet1_E   , Jet2_E    float64
	Jet1_EMf , Jet2_EMf  float64
	Jet1_Ntrk, Jet2_Ntrk int64
}
```

```go
// Array-typed tree input
type EventInArray struct {
	Jets_Px   [2]float64
	Jets_Py   [2]float64
	Jets_Pz   [2]float64
	Jets_E    [2]float64
	Jets_EMf  [2]float64
	Jets_Ntrk [2]int64
}
```