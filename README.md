# Simple Go examples

This repository contains few examples in go, mostly connected to high energy physics (HEP). Mostly:
  + `plotting`: very short example of plots produced in go
  + `CLs`: computation of a stat-only CLs exclusion [wikipedia](https://en.wikipedia.org/wiki/CLs_method_(particle_physics))
  + `reading-root-ttree`: event loop based on a `ROOT::TTree`, the main HEP software.


### CLs exclusion


### Reading a TTree

In this example, the initial TTree was produced from a LHE file [[arXiv:hep-ph/0609017](https://arxiv.org/abs/hep-ph/0609017)]
describing 10000 top-antitop quark productions at the LHE, as predicted by MadGraph
tool [[arXiv:1405.0301](https://arxiv.org/abs/1405.0301)], ran at the leading order.