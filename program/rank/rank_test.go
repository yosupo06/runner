package rank_test

import (
	"github.com/yosupo06/runner/program/rank"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestComment(te *testing.T) {
	var m sync.WaitGroup
	N := 100
	m.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer m.Done()
			s := strconv.Itoa(i)
			id := "test" + s
			for j := 0; j < 100; j++ {
				t := strconv.Itoa(j)
				co := "testcomment" + t
				rank.ChangeComment(id, co)
				r := rank.GetRanking()
				f := false
				for _, d := range r {
					if d.Id == id {
						if f {
							te.Fatal("Multi User")
						}
						if d.Comment != co {
							te.Error("Comment get error: ", d.Comment)
						} else {
							f = true
						}
					}
				}
				if !f {
					te.Error("comment get failed")
				}
			}
		}(i)
	}
	m.Wait()
}

func TestPoint(te *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var m sync.WaitGroup
	N := 100
	m.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer m.Done()
			s := strconv.Itoa(i)
			id := "test" + s
			sm := 0
			for j := 0; j < 100; j++ {
				d := rand.Intn(10000)
				sm += d
				rank.AddPoint(id, d)
				r := rank.GetRanking()
				f := false
				for _, d := range r {
					if d.Id == id {
						if f {
							te.Fatal("Multi User")
						}
						if d.Point != sm {
							te.Fatal("Error get point", d.Point)
						} else {
							f = true
						}
					}
				}
				if !f {
					te.Error("comment get failed")
				}
			}
		}(i)
	}
	m.Wait()
}
