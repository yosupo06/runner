package rank

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"sync"
)

type Data struct {
	Comment string
	Point   int
}

type RankData struct {
	Id string
	Data
}

var (
	rm   = new(sync.Mutex)
	rank = make(map[string]Data)
)

var session *mgo.Session

func init() {
	ses, err := mgo.Dial("localhost")
	session = ses
	if err != nil {
		panic(err)
	}
	ses.SetMode(mgo.Monotonic, true)
	c := ses.DB("runner").C("ranking")
	c.EnsureIndex(mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	})
	r := make([]RankData, 0)
	c.Find(nil).All(&r)
	for _, d := range r {
		rank[d.Id] = d.Data
	}
	go updateDB()
}

var (
	upMutex = new(sync.Mutex)
	upList  = make(chan string, 10000)
)

func updateDB() {
	ses := session.Copy()
	defer ses.Close()
	for id := range upList {
		c := ses.DB("runner").C("ranking")
		rm.Lock()
		c.Upsert(bson.M{"id": id}, RankData{id, rank[id]})
		rm.Unlock()
	}
}

func AddPoint(id string, p int) {
	rm.Lock()
	defer rm.Unlock()
	r := rank[id]
	r.Point += p
	rank[id] = r
	upMutex.Lock()
	upList <- id
	upMutex.Unlock()
}

func ChangeComment(id, c string) {
	rm.Lock()
	defer rm.Unlock()
	r := rank[id]
	r.Comment = c
	rank[id] = r
	upMutex.Lock()
	upList <- id
	upMutex.Unlock()
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
		r = append(r, RankData{k, d})
	}
	sort.Sort(sort.Reverse(ByPoint(r)))
	return r
}
