package internal

import (
	"errors"
	"math"
	"sort"

	"github.com/ArshiaDadras/Ariadne/pkg"
)

const (
	Sigma                = 5.0
	Beta                 = 10.0
	MaxDiffDistance      = 2000.0
	MaxCandidates        = 10
	MaxCandidateDistance = 200.0
	CandidateDistance    = 1.4142135624 * MaxCandidateDistance
)

var (
	ErrNoPathFound = errors.New("no path found")
)

func BestMatch(graph *pkg.Graph, points []GPSPoint) ([]*pkg.Edge, error) {
	if len(points) == 0 {
		return []*pkg.Edge{}, nil
	}
	dp, par := initializeDPAndPar(len(points))
	initializeValues(graph, points[0], dp)
	filterCandidates(dp[0])

	for i := 1; i < len(points); i++ {
		if len(dp[i-1]) == 0 {
			return nil, ErrNoPathFound
		}
		normalizeValues(dp[i-1])

		viterbi(graph, points, dp, par, i)
		filterCandidates(dp[i])
	}

	return bestPath(graph, points, dp, par)
}

func initializeDPAndPar(n int) ([]map[*pkg.Edge]float64, []map[*pkg.Edge]*pkg.Edge) {
	dp := make([]map[*pkg.Edge]float64, n)
	par := make([]map[*pkg.Edge]*pkg.Edge, n)
	for i := 0; i < n; i++ {
		dp[i] = make(map[*pkg.Edge]float64)
		par[i] = make(map[*pkg.Edge]*pkg.Edge)
	}
	return dp, par
}

func initializeValues(graph *pkg.Graph, initial GPSPoint, dp []map[*pkg.Edge]float64) {
	for _, candidate := range graph.Seg.Get(initial.Location, CandidateDistance) {
		dp[0][candidate] = EmmisionLogProbability(initial.Location.Distance(initial.Location.ClosestPointOnEdge(candidate)), Sigma)
	}
}

func normalizeValues(values map[*pkg.Edge]float64) {
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

func bestPath(graph *pkg.Graph, points []GPSPoint, dp []map[*pkg.Edge]float64, par []map[*pkg.Edge]*pkg.Edge) ([]*pkg.Edge, error) {
	best, edge := math.Inf(-1), (*pkg.Edge)(nil)
	for candidate, prob := range dp[len(points)-1] {
		if prob > best {
			best, edge = prob, candidate
		}
	}
	if edge == nil {
		return nil, ErrNoPathFound
	}

	result := make([]*pkg.Edge, 0)
	for i := len(points) - 1; i > 0; i-- {
		if par[i][edge].ID != edge.ID {
			path, err := graph.GetBestPath(graph.Nodes[edge.Start], graph.Nodes[par[i][edge].End], points[i].Location.Distance(points[i-1].Location)+MaxDiffDistance, true)
			if err != nil {
				return nil, err
			}

			result = append(result, edge)
			result = append(result, path...)
		}
		edge = par[i][edge]
	}
	return result, nil
}

func viterbi(graph *pkg.Graph, points []GPSPoint, dp []map[*pkg.Edge]float64, par []map[*pkg.Edge]*pkg.Edge, i int) {
	d1 := points[i].Location.Distance(points[i-1].Location)
	for _, candidate := range graph.Seg.Get(points[i].Location, CandidateDistance) {
		best, prv := math.Inf(-1), (*pkg.Edge)(nil)
		for prev, prevProb := range dp[i-1] {
			d2, err := roadDistance(graph, prev, candidate, points[i-1], points[i])
			if err != nil {
				continue
			}

			prob := prevProb + TransitionLogProbability(d1, d2, Beta)
			if prob > best {
				best, prv = prob, prev
			}
		}

		best += EmmisionLogProbability(points[i].Location.Distance(points[i].Location.ClosestPointOnEdge(candidate)), Sigma)
		if prv != nil {
			dp[i][candidate] = best
			par[i][candidate] = prv
		}
	}
}

func roadDistance(graph *pkg.Graph, prev, candidate *pkg.Edge, prevPoint, candidatePoint GPSPoint) (float64, error) {
	if prev.ID == candidate.ID {
		p1, p2 := prevPoint.Location.ClosestPointOnEdge(candidate), candidatePoint.Location.ClosestPointOnEdge(candidate)
		return candidate.LengthFrom(p1) - candidate.LengthFrom(p2), nil
	}
	d, err := graph.GetDistance(graph.Nodes[candidate.Start], graph.Nodes[prev.End], prevPoint.Location.Distance(candidatePoint.Location)+MaxDiffDistance, true)
	if err != nil {
		return 0, err
	}
	return d + prev.LengthFrom(prevPoint.Location.ClosestPointOnEdge(prev)) + candidate.LengthTo(candidatePoint.Location.ClosestPointOnEdge(candidate)), nil
}

func filterCandidates(values map[*pkg.Edge]float64) {
	if len(values) <= MaxCandidates {
		return
	}

	candidates := make([]*pkg.Edge, 0, len(values))
	for candidate := range values {
		candidates = append(candidates, candidate)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return values[candidates[i]] > values[candidates[j]]
	})

	for i := MaxCandidates; i < len(candidates); i++ {
		delete(values, candidates[i])
	}
}
