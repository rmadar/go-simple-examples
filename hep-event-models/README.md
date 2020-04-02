# Different event-models at once

Exemple of how to use go interface to define several HEP event models. The goal is to have a single event model at the user level, and handle the differences & conversion in thanks to an interface. This example is based on ROOT object translated in go in groot (go-hep package), but can be easily generalized.

The use case in mind for this example is to have two inputs event formats, `EventInFlat` and `EventInArrays`, for one unique output format `EventOut`. The user just has to define **once**, in an interface, the format of `EventInFlat` and `EventInFlat`, and how to translate each of them into the format `EventOut`. Then, with a simple flag (passed as argument) one can run the same executable on the two input formats

In summary, we the following workflow:
```
| EventInFlat reading  |
|                      | --> EventOut --> Event processing --> Output writing
| EventInArray reading |
```
  1. EventOut, to be processed (add new variables) and saved
  1. EventInFlat, to be loaded and converted into EventOut
  1. EventInArray, to be loaded and converted into EventOut
where the two last are satisfying the single interface 'EventIn'.

For this example, we consider events with 2-jets in the final state, described by 12 numbers, 2 x (4-vectors + number of tracks + EM fraction), that can be organized in 3 different ways.

