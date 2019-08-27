package world

// MESSAGE
//
// PLAYER
// DIR 위 : 0, 오른쪽 1, 아래 : 2, 왼쪽 : 3

// CS
// Owner 0 : 중립 , 1 ~ 플레이어 id

type Message struct{
	Msg string `json:"msg,omitempty"`
	Type int `json:"type,omitempty"` //ENTER
	Player Player `json:"player,omitempty"`
	Score map[string]int `json:"score,omitempty"`
	Players map[string]*Player `json:"players,omitempty"`
}
