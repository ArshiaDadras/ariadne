package internal

import "math"

func EmmisionLogProbability(x, sigma float64) float64 {
	return -0.5 * (x * x) / (sigma * sigma)
}

func EmissionProbability(x, sigma float64) float64 {
	return math.Exp(-0.5*(x*x)/(sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
}

func TransitionLogProbability(x, y, beta float64) float64 {
	return -math.Abs(x-y) / beta
}

func TransitionProbability(x, y, beta float64) float64 {
	return math.Exp(-math.Abs(x-y)/beta) / beta
}
