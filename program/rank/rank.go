package rank

import (
	"fmt"
	"sort"
	"sync"
)

type Data struct {
	Comment string
	Point   int
}

var (
	rm   = new(sync.Mutex)
	rank = make(map[string]Data)
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
	fmt.Println("GetRankData ", rank)
	r := make([]RankData, 0, len(rank))
	for k, d := range rank {
		fmt.Println("get", k, d)
		r = append(r, RankData{k, d.Comment, d.Point})
	}
	rm.Unlock()
	fmt.Println("Data", r)
	sort.Sort(sort.Reverse(ByPoint(r)))
	return r
}
