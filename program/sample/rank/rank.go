package rank

import (
	"sort"
	"sync"
)

var (
	rm   = new(sync.Mutex)
	rank = make(map[string]struct {
		Comment string
		Point   int
	})
)

func AddPoint(id string, p int) {
	rm.Lock()
	defer rm.Unlock()
	r := rank[id]
	r.Point += p
	rank[id] = r
}

func ChangeComment(id, c string) {
	rm.Lock()
	defer rm.Unlock()
	r := rank[id]
	r.Comment = c
	rank[id] = r
}

type RankData struct {
	Id      string
	Comment string
	Point   int
}

type ByPoint []RankData

func (a ByPoint) Len() int           { return len(a) }
func (a ByPoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPoint) Less(i, j int) bool { return a[i].Point < a[j].Point }

func GetRanking() []RankData {
	rm.Lock()
	defer rm.Unlock()
	r := make([]RankData, 0, len(rank))
	for k, d := range rank {
		r = append(r, RankData{k, d.Comment, d.Point})
	}
	sort.Sort(sort.Reverse(ByPoint(r)))
	return r
}
