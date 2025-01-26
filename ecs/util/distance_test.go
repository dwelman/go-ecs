package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/gonum/spatial/r3"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		name           string
		v1             interface{}
		v2             interface{}
		expectedResult float64
		expectedError  error
	}{
		{
			name:           "Distance2D - no change",
			v1:             r2.Vec{X: 1, Y: 1},
			v2:             r2.Vec{X: 1, Y: 1},
			expectedResult: 0,
			expectedError:  nil,
		},
		{
			name:           "Distance3D - no change",
			v1:             r3.Vec{X: 1, Y: 1, Z: 1},
			v2:             r3.Vec{X: 1, Y: 1, Z: 1},
			expectedResult: 0,
			expectedError:  nil,
		},
		{
			name:           "Distance2D - different X",
			v1:             r2.Vec{X: 1, Y: 1},
			v2:             r2.Vec{X: 2, Y: 1},
			expectedResult: 1,
			expectedError:  nil,
		},
		{
			name:           "Distance3D - different X",
			v1:             r3.Vec{X: 1, Y: 1, Z: 1},
			v2:             r3.Vec{X: 2, Y: 1, Z: 1},
			expectedResult: 1,
			expectedError:  nil,
		},
		{
			name:           "Distance2D - different Y",
			v1:             r2.Vec{X: 1, Y: 1},
			v2:             r2.Vec{X: 1, Y: 2},
			expectedResult: 1,
			expectedError:  nil,
		},
		{
			name:           "Distance3D - different Y",
			v1:             r3.Vec{X: 1, Y: 1, Z: 1},
			v2:             r3.Vec{X: 1, Y: 2, Z: 1},
			expectedResult: 1,
			expectedError:  nil,
		},
		{
			name:           "Distance2D - different X, Y",
			v1:             r2.Vec{X: 1, Y: 1},
			v2:             r2.Vec{X: 2, Y: 2},
			expectedResult: 1.4142135623730951,
			expectedError:  nil,
		},
		{
			name:           "Distance3D - different X, Y, Z",
			v1:             r3.Vec{X: 1, Y: 1, Z: 1},
			v2:             r3.Vec{X: 2, Y: 2, Z: 2},
			expectedResult: 1.7320508075688772,
			expectedError:  nil,
		},
		{
			name:           "Distance - different types",
			v1:             r2.Vec{X: 1, Y: 1},
			v2:             r3.Vec{X: 1, Y: 1, Z: 1},
			expectedResult: 0,
			expectedError:  ErrMixedVectors,
		},
		{
			name:           "Distance - large difference",
			v1:             r3.Vec{X: 1, Y: 1, Z: 1},
			v2:             r3.Vec{X: 100, Y: 100, Z: 100},
			expectedResult: 171.47302994931886,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Distance(tt.v1, tt.v2)
			require.Equal(t, tt.expectedError, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
