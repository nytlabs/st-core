package core

import (
	"math"
	"math/rand"
)

// UniformRandom emits a uniform random between 0 and 1
func UniformRandom() Spec {
	return Spec{
		Name:    "uniform",
		Inputs:  []Pin{},
		Outputs: []Pin{Pin{"draw", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			out[0] = rand.Float64()
			return nil
		},
	}
}

// NormalRandom emits a normally distributed random number with the
// supplied mean and variance
func NormalRandom() Spec {
	return Spec{
		Name:    "normal",
		Inputs:  []Pin{Pin{"mean", NUMBER}, Pin{"variance", NUMBER}},
		Outputs: []Pin{Pin{"draw", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			variance, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("variance must be a number")
				return nil
			}
			mean, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("mean must be a number")
				return nil
			}
			out[0] = rand.NormFloat64()*math.Sqrt(variance) + mean
			return nil
		},
	}
}

// the global random number source
var RAND *rand.Rand = rand.New(rand.NewSource(12345))

// ZipfRandom emits a Zipfian distributed random number
// notation follows the wikipedia page http://en.wikipedia.org/wiki/Zipf%E2%80%93Mandelbrot_law not the golang Zipf parameters
func ZipfRandom() Spec {
	return Spec{
		Name: "Zipf",
		Inputs: []Pin{
			Pin{"q", NUMBER}, Pin{"s", NUMBER}, Pin{"N", NUMBER}},
		Outputs: []Pin{Pin{"draw", NUMBER}},
		Kernel: func(in, out, internal MessageMap, ss Source, i chan Interrupt) Interrupt {

			q, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("q must be a number")
				return nil
			}
			s, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("s must be a number")
				return nil
			}
			N, ok := in[2].(float64)
			if !ok {
				out[0] = NewError("N must be an number")
				return nil
			}

			z := rand.NewZipf(RAND, s, q, uint64(N))
			out[0] = z.Uint64()
			return nil
		},
	}
}

// poisson returns an integer (though we actually pretend it's a float) from a Poisson distrbution
func poisson(λ float64) float64 {
	var k float64
	L := math.Exp(-λ)
	k = 0
	p := 1.0
	for {
		k++
		u := RAND.Float64()
		p = p * u
		if p <= L {
			return k - 1
		}
	}
}

// PoissonRandom emits a Poisson distribtued random number
func PoissonRandom() Spec {
	return Spec{
		Name:    "poisson",
		Inputs:  []Pin{Pin{"rate", NUMBER}},
		Outputs: []Pin{Pin{"draw", NUMBER}},
		Kernel: func(in, out, internal MessageMap, ss Source, i chan Interrupt) Interrupt {
			λ, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("rate must be a number")
				return nil
			}
			if λ < 0 {
				out[0] = NewError("rate must be positive")
				return nil
			}
			out[0] = poisson(λ)
			return nil
		},
	}
}

// ExponentialRandom emits an Exponentially distribtued random number
func ExponentialRandom() Spec {
	return Spec{
		Name:    "exponential",
		Inputs:  []Pin{Pin{"rate", NUMBER}},
		Outputs: []Pin{Pin{"draw", NUMBER}},
		Kernel: func(in, out, internal MessageMap, ss Source, i chan Interrupt) Interrupt {
			λ, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("rate must be a number")
				return nil
			}
			if λ < 0 {
				out[0] = NewError("rate must be positive")
				return nil
			}
			out[0] = rand.ExpFloat64() / λ
			return nil
		},
	}
}

// BernoulliRandom emits a draw from a Bernoulli distribution. This block returns a boolean
func BernoulliRandom() Spec {
	return Spec{
		Name:    "bernoulli",
		Inputs:  []Pin{Pin{"bias", NUMBER}},
		Outputs: []Pin{Pin{"draw", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			r := RAND.Float64()
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("bias must be a number")
				return nil
			}
			if p < 0 || p > 1 {
				out[0] = NewError("bias must be between 0 and 1")
			}
			if r > p {
				out[0] = false
			} else {
				out[0] = true
			}
			return nil
		},
	}
}
