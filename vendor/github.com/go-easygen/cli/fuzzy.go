package cli

func match(s, t string) (float32, bool) {
	return matchWithMinDifferRate(s, t, 0.3)
}

func matchWithMinDifferRate(s, t string, minDifferRate float32) (float32, bool) {
	dist := editDistance([]byte(s), []byte(t))
	differRate := float32(dist) / float32(max(len(s), len(t))+4)
	return differRate, differRate <= minDifferRate
}

func editDistance(s, t []byte) float32 {
	var (
		m = len(s)
		n = len(t)
		d = make([][]float32, m+1)
	)
	for i := 0; i < m+1; i++ {
		d[i] = make([]float32, n+1)
		d[i][0] = float32(i)
	}
	for j := 0; j < n+1; j++ {
		d[0][j] = float32(j)
	}

	for j := 1; j < n+1; j++ {
		for i := 1; i < m+1; i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+1))
			}
		}
	}

	return d[m][n]
}

func min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type editDistanceRank struct {
	s string
	d float32
}

type editDistanceRankSlice []editDistanceRank

func (dists editDistanceRankSlice) Len() int {
	return len(dists)
}

func (dists editDistanceRankSlice) Less(i, j int) bool {
	return dists[i].d < dists[j].d
}

func (dists editDistanceRankSlice) Swap(i, j int) {
	dists[i], dists[j] = dists[j], dists[i]
}
