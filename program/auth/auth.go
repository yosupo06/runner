package auth

import (
	"crypto/rand"
	"errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"net/http"
	"strconv"
)

type User struct {
	Id    string
	Token string
	Pass  string //must be salted
}

func init() {
	ses, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer ses.Close()
	c := ses.DB("runner").C("user")
	c.EnsureIndex(mgo.Index{
		Key:    []string{"id"},
		Unique: true,
	})
}

func makeToken() string {
	c := 36
	b := make([]byte, c)
	rand.Read(b)
	r := ""
	for i := 0; i < c; i++ {
		d, _ := rand.Int(rand.Reader, big.NewInt(36))
		r += strconv.FormatInt(d.Int64(), 36)
	}
	return r
}

func AddUser(id string, pass string) error {
	ses, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer ses.Close()
	c := ses.DB("runner").C("user")
	co, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		panic(err)
	}
	if co != 0 {
		return errors.New("このIDはもう使われています")
	}
	c.Insert(User{id, makeToken(), pass})
	return nil
}

func GetUser(id string) (*User, bool) {
	ses, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer ses.Close()

	c := ses.DB("runner").C("user")
	var u User
	err = c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return nil, false
	}
	return &u, true
}

func AuthPass(id string, pass string) bool {
	u, ok := GetUser(id)
	if !ok {
		return false
	}
	return u.Pass == pass
}

func AuthToken(id string, token string) bool {
	u, ok := GetUser(id)
	if !ok {
		return false
	}
	if u.Token != token {
		return false
	}
	return true
}

func SetCookie(rw http.ResponseWriter, id string) {
	http.SetCookie(rw,
		&http.Cookie{
			Name:   "id",
			Value:  id,
			MaxAge: 60 * 60 * 24 * 14,
		})
}

func DelCookie(rw http.ResponseWriter) {
	http.SetCookie(rw,
		&http.Cookie{
			Name:   "id",
			MaxAge: -1,
		})
}

func GetCookie(req *http.Request) (*User, bool) {
	id, err := req.Cookie("id")
	if err != nil {
		return nil, false
	}
	return GetUser(id.Value)
}
