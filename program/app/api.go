package app

import (
	"fmt"
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/rank"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type TL struct {
	Vote time.Time //sampling
	Comm time.Time //comment
	Hist time.Time //submit
}

type H struct {
	Rabit  int
	Points [5]int
}

var (
	hm   = new(sync.RWMutex)
	hist = make([]H, 0)
)

var (
	pm    = new(sync.Mutex) //price and mp
	price = make(map[string]int)
	mp    = 100
)

var votec = 0

var (
	tm = new(sync.Mutex) //time limit
	tl = make(map[string]TL)
)

func owner() {
	c := time.Tick(10 * time.Second)
	for range c {
		if time.Now().Before(Start) {
			continue
		}
		pm.Lock()
		p := make([]int, 0)
		for _, d := range price {
			p = append(p, d)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(p)))
		h := H{}
		h.Rabit = rand.Intn(mp) + 1

		ma := 0
		if len(p) != 0 {
			ma = p[0]
		}
		if ma < h.Rabit {
			ma = h.Rabit
		}

		le := len(p)
		if 5 < le {
			le = 5
		}
		for i := 0; i < le; i++ {
			h.Points[i] = p[i]
		}

		hm.Lock()
		hist = append(hist, h)
		hm.Unlock()

		type P struct {
			Id string
			P  int
		}
		u := make([]P, 0)
		for k, d := range price {
			if d == ma {
				continue
			}
			u = append(u, P{k, d})
		}
		price = make(map[string]int)
		votec++
		if votec%100 == 0 {
			mp += 100
		}

		pm.Unlock()

		for _, d := range u {
			rank.AddPoint(d.Id, d.P)
		}
	}
}

func init() {
	go owner()
}

func getForm(req *http.Request) (string, string) {
	return req.FormValue("id"),
		req.FormValue("token")
}

func VoteApi(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(Start) {
		NotFound(rw, req)
		return
	}
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		fmt.Fprintln(rw, "Invalid Token")
		return
	}
	n := time.Now()
	tm.Lock()
	t := tl[id]
	if n.Before(t.Vote) {
		fmt.Fprintln(rw, "TLE")
		return
	}
	u, err := strconv.Atoi(req.FormValue("price"))
	if err != nil {
		fmt.Fprintln(rw, "Format Error")
		return
	}
	pm.Lock()
	if u <= 0 || mp < u {
		fmt.Fprintf(rw, "price must be in range [0, %d]", mp)
		return
	}
	price[id] = u
	pm.Unlock()
	t.Vote = n.Add(time.Second)
	tl[id] = t
	tm.Unlock()
}

func HistoryApi(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(Start) {
		NotFound(rw, req)
		return
	}
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		fmt.Fprintln(rw, "Invalid Token")
		return
	}
	tm.Lock()
	n := time.Now()
	t := tl[id]
	if n.Before(t.Hist) {
		fmt.Fprintln(rw, "TLE")
	}
	rank.ChangeComment(id, req.FormValue("comment"))
	t.Hist = n.Add(time.Second)
	tl[id] = t
	tm.Unlock()
	hm.RLock()
	hl := len(hist) - 5
	if hl < 0 {
		hl = 0
	}
	for _, d := range hist[hl:] {
		fmt.Fprintln(rw, d.Rabit, d.Points)
	}
	hm.RUnlock()
}

func GetHistory() []H {
	hm.RLock()
	fmt.Println("Hist ", hist)
	h := make([]H, len(hist))
	copy(h, hist)
	hm.RUnlock()
	return h
}

func CommentApi(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(Start) {
		NotFound(rw, req)
		return
	}
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		fmt.Fprintln(rw, "Invalid Token")
		return
	}
	tm.Lock()
	n := time.Now()
	t := tl[id]
	if n.Before(t.Comm) {
		fmt.Fprintln(rw, "TLE")
	}
	rank.ChangeComment(id, req.FormValue("comment"))
	t.Comm = n.Add(time.Second)
	tl[id] = t
	tm.Unlock()
}
