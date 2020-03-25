module github.com/rmadar/go-simple-examples

go 1.14

require (
	github.com/golang/geo v0.0.0-20200319012246-673a6f80352d
	github.com/rmadar/go-lorentz-vector v0.0.0-20200325122951-f62ad51a1848
	go-hep.org/x/hep v0.24.1
	gonum.org/v1/gonum v0.7.0
	gonum.org/v1/plot v0.7.0
)

// For local tests during dev 
//replace github.com/rmadar/go-lorentz-vector => /home/rmadar/cernbox/goDev/go-lorentz-vector
