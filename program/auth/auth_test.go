package auth_test

import (
	. "github.com/yosupo06/runner/program/auth"
	"strconv"
	"testing"
)

func TestUser(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			s := strconv.Itoa(i)
			err := AddUser("test"+s, "test"+s)
			if err != nil {
				t.Fatal(err)
			}
			u, ok := GetUser("test" + s)
			if ok {
				t.Fatal("get failed")
			}
			if u.Id != "test"+s {
				t.Fatal("get name error")
			}
			if AuthToken("test"+s, u.Token) {
				t.Fatal("error auth token")
			}
			if AuthPass("test"+s, "test"+s) {
				t.Fatal("error auth pass")
			}
		}()
	}
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
