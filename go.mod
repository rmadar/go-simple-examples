module github.com/rmadar/go-simple-examples

go 1.14

require (
	github.com/golang/geo v0.0.0-20200319012246-673a6f80352d
	github.com/rmadar/go-lorentz-vector v0.0.0-20200327223025-1ae75e287f6e
	go-hep.org/x/hep v0.27.1-0.20200605151707-45f0797ac50e
	golang.org/x/exp v0.0.0-20200331195152-e8c3332aa8e5
	golang.org/x/mobile v0.0.0-20200329125638-4c31acba0007 // indirect
	gonum.org/v1/gonum v0.7.1-0.20200602002949-4ff1bb0a480e
	gonum.org/v1/plot v0.7.1-0.20200602093449-6d232e045386
)

// For local tests during dev
// replace github.com/rmadar/go-lorentz-vector => /home/rmadar/cernbox/goDev/go-lorentz-vector
