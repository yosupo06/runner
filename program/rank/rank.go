package rank

import ()

type RankData struct {
	id      string
	comment string
	point   float64
}

var rank = make([]RankData, 0)

func AddPoint() {

}

func ChangeComment(id, c string) {
	for i, _ := range rank {
		if rank[i].id == id {
			rank[i].comment = c
			return
		}
	}
	rank = append(rank, RankData{id, c, 0.0})
}

func GetRanking() []RankData {
	return rank
}
