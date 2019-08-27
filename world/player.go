package world

type Player struct{
	ID string `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	DIR int `json:"dir"`
	CS []CS `json:"cs"`
	GetCS int `json:"getCs"`
	Score int `json:"score"`
	Team int `json:"team"`
}

func NewPlayer(id string) *Player{
	p := Player{}
	p.ID = id
	p.X = 10
	p.Y = 210
	p.DIR = 0
	p.Score = 0
	p.Team = 0
	p.GetCS = 0
	return &p
}