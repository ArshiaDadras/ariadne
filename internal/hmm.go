package internal

import (
	"errors"
	"math"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

const (
	Sigma                = 5.0
	Beta                 = 10.0
	MaxDiffDistance      = 2000.0
	MaxCandidateDistance = 200.0
)

var (
	ErrNoPathFound = errors.New("no path found")
)

func BestMatch(graph *pkg.Graph, points []GPSPoint) ([]*pkg.Edge, error) {
	dp, par := initializeDPAndPar(len(points))
	initializeValues(graph, points, dp)

	for i := 1; i < len(points); i++ {
		normalizeValues(dp[i-1])

		prevNodes := make([]*pkg.Node, 0, len(dp[i-1]))
		for key := range dp[i-1] {
			prevNodes = append(prevNodes, key)
		}
		if len(prevNodes) == 0 {
			return nil, ErrNoPathFound
		}

		viterbi(graph, points, dp, par, i, prevNodes)
	}

	return bestPath(graph, points, dp, par)
}

func initializeDPAndPar(n int) ([]map[*pkg.Node]float64, []map[*pkg.Node]*pkg.Node) {
	dp := make([]map[*pkg.Node]float64, n)
	par := make([]map[*pkg.Node]*pkg.Node, n)
	for i := 0; i < n; i++ {
		dp[i] = make(map[*pkg.Node]float64)
		par[i] = make(map[*pkg.Node]*pkg.Node)
	}
	return dp, par
}

func initializeValues(graph *pkg.Graph, points []GPSPoint, dp []map[*pkg.Node]float64) {
	candidates := graph.GetSquare(points[0].Location, MaxCandidateDistance)
	for _, candidate := range candidates {
		dp[0][candidate] = EmmisionLogProbability(points[0].Location.Distance(candidate.Position), Sigma)
	}
}

func normalizeValues(values map[*pkg.Node]float64) {
	best := math.Inf(-1)
	for _, prob := range values {
		if prob > best {
			best = prob
		}
	}
	for key := range values {
		values[key] -= best
	}
}

func bestPath(graph *pkg.Graph, points []GPSPoint, dp []map[*pkg.Node]float64, par []map[*pkg.Node]*pkg.Node) ([]*pkg.Edge, error) {
	best, node := math.Inf(-1), (*pkg.Node)(nil)
	for candidate, prob := range dp[len(points)-1] {
		if prob > best {
			best, node = prob, candidate
		}
	}
	if node == nil {
		return nil, ErrNoPathFound
	}

	result := make([]*pkg.Edge, 0)
	for i := len(points) - 1; i > 0; i-- {
		if par[i][node].ID != node.ID {
			path, err := graph.GetBestPath(node, par[i][node], points[i].Location.Distance(points[i-1].Location)+MaxDiffDistance, true)
			if err != nil {
				return nil, err
			}
			result = append(result, path...)
		}
		node = par[i][node]
	}
	return result, nil
}

func viterbi(graph *pkg.Graph, points []GPSPoint, dp []map[*pkg.Node]float64, par []map[*pkg.Node]*pkg.Node, i int, prevNodes []*pkg.Node) {
	d1 := points[i].Location.Distance(points[i-1].Location)
	for _, candidate := range graph.GetSquare(points[i].Location, MaxCandidateDistance) {
		best, prv := math.Inf(-1), (*pkg.Node)(nil)
		for _, prev := range prevNodes {
			d2, err := graph.GetDistance(candidate, prev, d1+MaxDiffDistance, true)
			if err != nil {
				continue
			}

			prob := dp[i-1][prev] + TransitionLogProbability(d1, d2, Beta)
			if prob > best {
				best, prv = prob, prev
			}
		}

		best += EmmisionLogProbability(points[i].Location.Distance(candidate.Position), Sigma)
		if prv != nil {
			dp[i][candidate] = best
			par[i][candidate] = prv
		}
	}
}
