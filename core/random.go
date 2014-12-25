package core

import (
	"math"
	"math/rand"
)

// UniformRandom emits a uniform random between 0 and 1
func UniformRandom() Spec {
	return Spec{
		Inputs:  []Pin{},
		Outputs: []Pin{Pin{"draw"}},
		Kernel: func(in, out MessageMap, s StateLocker, i chan Interrupt) Interrupt {
			out[0] = rand.Float64()
			return nil
		},
	}
}

// NormalRandom emits a normally distributed random number with the
// supplied mean and variance
func NormalRandom() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"mean"}, Pin{"variance"}},
		Outputs: []Pin{Pin{"draw"}},
		Kernel: func(in, out MessageMap, s StateLocker, i chan Interrupt) Interrupt {
			variance, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("variance must be a float")
				return nil
			}
			mean, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("mean must be a float")
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
		Inputs:  []Pin{Pin{"q"}, Pin{"s"}, Pin{"N"}},
		Outputs: []Pin{Pin{"draw"}},
		Kernel: func(in, out MessageMap, ss StateLocker, i chan Interrupt) Interrupt {

			q, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("q must be a float")
				return nil
			}
			s, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("s must be a float")
				return nil
			}
			N, ok := in[2].(int)
			if !ok {
				out[0] = NewError("N must be an int")
				return nil
			}

			z := rand.NewZipf(RAND, s, q, uint64(N))
			out[0] = z.Uint64()
			return nil
		},
	}
}

func poisson(λ float64) int {
	L := math.Exp(-λ)
	k := 0
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
		Inputs:  []Pin{Pin{"rate"}},
		Outputs: []Pin{Pin{"draw"}},
		Kernel: func(in, out MessageMap, ss StateLocker, i chan Interrupt) Interrupt {
			λ, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("rate must be a float")
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

// BernoulliRandom emits a draw from a Bernoulli distribution. This block returns a boolean
func BernoulliRandom() Spec {
	return Spec{
		Inputs:  []Pin{Pin{"bias"}},
		Outputs: []Pin{Pin{"draw"}},
		Kernel: func(in, out MessageMap, ss StateLocker, i chan Interrupt) Interrupt {
			r := RAND.Float64()
			p, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("bias must be a float")
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
