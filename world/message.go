package world

// MESSAGE
//
// PLAYER
// DIR 위 : 0, 오른쪽 1, 아래 : 2, 왼쪽 : 3

// CS
// Owner 0 : 중립 , 1 ~ 플레이어 id

type CS struct{
	X int `json:"x"`
	Y int `json:"y"`
	Owner int `json:"owner""`
}
type Player struct{
	ID int `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	DIR int `json:"dir"`
	CS []CS `json:"cs"`
	Score int `json:"score"`
}

type Message struct{
	Msg string `json:"msg"`
	Player map[string]Player `json:"Player"`
	Type string `json:"type"`
	CS []CS `json:"CS"`
	Players []string
}

