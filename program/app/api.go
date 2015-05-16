package app

import (
	"fmt"
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/config"
	"github.com/yosupo06/runner/program/rank"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

type TL struct {
	Vote time.Time //sampling
	Comm time.Time //comment
	Rank time.Time //submit
	Info time.Time //info
}

var (
	pm     = new(sync.Mutex) //price and mp
	price  = make(map[string]int)
	mp     = 100
	votec  = 0
	rabNum = 0
	lastB  = 0
)

var (
	tm = new(sync.Mutex) //time limit
	tl = make(map[string]TL)
)

func owner() {
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := time.Tick(10 * time.Second)
	for range c {
		if time.Now().Before(config.Start) {
			continue
		}
		pm.Lock()
		p := make([]int, 0)
		for _, d := range price {
			p = append(p, d)
		}
		for i := 0; i < rabNum; i++ {
			p = append(p, gen.Intn(mp)+1)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(p)))

		rabNum = gen.Intn(len(price) + 1)
		lastB = 0
		if len(p) != 0 {
			lastB = p[len(p)/2] - 1
		}

		type P struct {
			Id string
			P  int
		}
		u := make([]P, 0)
		for k, d := range price {
			if d > lastB {
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

func errorApi(rw http.ResponseWriter, req *http.Request, mes string) {
	fmt.Fprintln(rw, "Error")
	fmt.Fprintln(rw, mes)
}

func VoteApi(rw http.ResponseWriter, req *http.Request) {
	n := time.Now()
	if n.Before(config.Start) {
		NotFound(rw, req)
		return
	}
	if n.After(config.End) {
		errorApi(rw, req, "Contest ended")
		return
	}
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		errorApi(rw, req, "Invalid Token")
		return
	}
	u, err := strconv.Atoi(req.FormValue("price"))
	if err != nil {
		errorApi(rw, req, "Format Error(Price)")
		return
	}

	tm.Lock()
	defer tm.Unlock()
	t := tl[id]
	if n.Before(t.Vote) {
		errorApi(rw, req, "TLE")
		return
	}

	pm.Lock()
	defer pm.Unlock()
	if u <= 0 || mp < u {
		errorApi(rw, req, "Range Error(Price)")
		return
	}
	price[id] = u
	t.Vote = n.Add(time.Second)
	tl[id] = t
	fmt.Fprintln(rw, "Success")
}

func CommentApi(rw http.ResponseWriter, req *http.Request) {
	n := time.Now()
	if n.Before(config.Start) {
		NotFound(rw, req)
		return
	}
	if n.After(config.End) {
		errorApi(rw, req, "Contest ended")
		return
	}
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		errorApi(rw, req, "Invalid Token")
		return
	}
	c := req.FormValue("comment")
	const ML = 1000
	if len(c) > ML {
		errorApi(rw, req, "Text too large")
		return
	}
	if !utf8.ValidString(c) {
		errorApi(rw, req, "Comment must be UTF-8")
		return
	}
	tm.Lock()
	defer tm.Unlock()
	t := tl[id]
	if n.Before(t.Comm) {
		errorApi(rw, req, "TLE")
		return
	}
	rank.ChangeComment(id, req.FormValue("comment"))
	t.Comm = n.Add(time.Second)
	tl[id] = t
	fmt.Fprintln(rw, "Success")
}

func InfoApi(rw http.ResponseWriter, req *http.Request) {
	n := time.Now()
	if n.Before(config.Start) {
		NotFound(rw, req)
		return
	}
	if n.After(config.End) {
		errorApi(rw, req, "Contest ended")
		return
	}

	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		errorApi(rw, req, "Invalid Token")
		return
	}

	tm.Lock()
	defer tm.Unlock()
	t := tl[id]
	if n.Before(t.Info) {
		errorApi(rw, req, "TLE")
		return
	}
	pm.Lock()
	fmt.Fprintln(rw, "Success")
	fmt.Fprintln(rw, votec)
	fmt.Fprintln(rw, lastB)
	fmt.Fprintln(rw, rabNum)
	pm.Unlock()
	t.Info = n.Add(time.Second)
	tl[id] = t

}

func RankingApi(rw http.ResponseWriter, req *http.Request) {
	n := time.Now()
	if n.Before(config.Start) {
		NotFound(rw, req)
		return
	}
	if n.After(config.End) {
		errorApi(rw, req, "Contest ended")
		return
	}

	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		errorApi(rw, req, "Invalid Token")
		return
	}

	tm.Lock()
	defer tm.Unlock()
	t := tl[id]
	if n.Before(t.Rank) {
		errorApi(rw, req, "TLE")
		return
	}
	r := rank.GetRanking()
	fmt.Fprintln(rw, "Success")
	fmt.Fprintln(rw, len(r))
	for _, d := range r {
		fmt.Fprintln(rw, d.Point, d.Id)
	}
	t.Rank = n.Add(time.Second)
	tl[id] = t
}
