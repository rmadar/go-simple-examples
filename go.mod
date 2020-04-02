module github.com/rmadar/go-simple-examples

go 1.14

require (
	github.com/golang/geo v0.0.0-20200319012246-673a6f80352d
	github.com/rmadar/go-lorentz-vector v0.0.0-20200327223025-1ae75e287f6e
	go-hep.org/x/hep v0.25.1-0.20200330174027-9e59fb905679
	gonum.org/v1/gonum v0.7.1-0.20200330111830-e98ce15ff236
	gonum.org/v1/plot v0.7.1-0.20200323092842-6973214b8663
)

// For local tests during dev
// replace github.com/rmadar/go-lorentz-vector => /home/rmadar/cernbox/goDev/go-lorentz-vector
