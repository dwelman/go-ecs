package util

import (
	"errors"

	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/gonum/spatial/r3"
)

var ErrMixedVectors = errors.New("cannot mix 2D and 3D vectors")

func Distance2D(v1 r2.Vec, v2 r2.Vec) float64 {
	return r2.Norm(r2.Sub(v1, v2))
}

func Distance3D(v1 r3.Vec, v2 r3.Vec) float64 {
	return r3.Norm(r3.Sub(v1, v2))
}

func Distance(v1 interface{}, v2 interface{}) (float64, error) {
	switch v1 := v1.(type) {
	case r2.Vec:
		if v2, ok := v2.(r2.Vec); ok {
			return Distance2D(v1, v2), nil
		}
	case r3.Vec:
		if v2, ok := v2.(r3.Vec); ok {
			return Distance3D(v1, v2), nil
		}
	}
	return 0, ErrMixedVectors
}
