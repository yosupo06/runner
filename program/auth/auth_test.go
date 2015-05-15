package auth_test

import (
	. "github.com/yosupo06/runner/program/auth"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var m sync.WaitGroup
	m.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer m.Done()
			s := strconv.Itoa(rand.Int())
			id := "test" + s
			err := AddUser(id, id)
			if err != nil {
				t.Fatal(err)
			}
			u, ok := GetUser(id)
			if !ok {
				t.Fatal("get failed")
			}
			if u.Id != id {
				t.Fatal("get name error")
			}
			if !AuthToken(id, u.Token) {
				t.Fatal("error auth token")
			}
			if !AuthPass(id, id) {
				t.Fatal("error auth pass")
			}
		}(i)
	}
	m.Wait()
}

func BenchmarkToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeToken()
	}
}

func BenchmarkHash(b *testing.B) {
	s := "test~i'm morita ~ LALALA~"
	for i := 0; i < b.N; i++ {
		Hash(s)
	}
}
