package downsample

import (
	"math"
	"sort"
)

func NewPoints() Points {
	return make(Points, 0)
}

func (ps Points) Downsample(newStep int) Points {
	if len(ps) == 0 {
		return []Point{}
	}

	if !sort.IsSorted(ps) {
		sort.Sort(ps)
	}

	if len(ps) == 1 {
		return []Point{
			Point{ps[0].Timestamp, ps[0].Value},
		}
	}

	start := ps[0].Timestamp
	end := ps[len(ps)-1].Timestamp

	// 不应该出现这种case
	if start > end {
		return []Point{}
	}

	// 不应该出现这种case
	if start == end {
		return []Point{
			Point{ps[0].Timestamp, ps[0].Value},
		}
	}

	res := make([]Point, 0)
	for t1 := start; t1 <= end; t1 += int64(newStep) {
		r := ps.avg(t1, t1+int64(newStep), end)
		res = append(res, Point{t1, r})
	}

	return res
}

func (ps Points) avg(t1, t2 int64, max int64) float64 {
	l := len(ps)
	var (
		count int = 0
		sum   float64
	)
	idx := sort.Search(l, func(i int) bool { return ps[i].Timestamp >= t1 })
	if idx == l {
		return math.NaN()
	}

	if t2 > max {
		for _, p := range ps {
			if p.Timestamp >= ps[idx].Timestamp {
				if v := float64(p.Value); !math.IsNaN(v) {
					sum += v
					count++
				}
			}
		}
	} else {
		for _, p := range ps {
			if p.Timestamp >= ps[idx].Timestamp && p.Timestamp < t2 {
				if v := float64(p.Value); !math.IsNaN(v) {
					sum += v
					count++
				}
			}
		}
	}

	if count == 0 {
		return math.NaN()
	}

	return sum / float64(count)
}

type Points []Point

type Point struct {
	Timestamp int64
	Value     float64
}

func (ps Points) Less(i, j int) bool { return ps[i].Timestamp < ps[j].Timestamp }
func (ps Points) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps Points) Len() int           { return len(ps) }
