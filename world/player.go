package world

type Player struct{
	ID string `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	DIR int `json:"dir"`
	CS []CS `json:"cs"`
	Score int `json:"score"`
	Team int `json:"team"`
}

func NewPlayer(id string) *Player{
	p := &Player{}
	p.ID = id
	p.X = 0
	p.Y = 0
	p.DIR = 0
	p.Score = 0
	p.Team = 0
	p.CS = []CS{}
	return p
}