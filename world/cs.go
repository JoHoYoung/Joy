package world

type CS struct{
	X int `json:"x"`
	Y int `json:"y"`
	Owner int `json:"owner"`
}

func newCS() *CS{
	s := &CS{}
	s.X = 10
	s.Y = 10
	return s
}


func collision(player Player, cs CS) bool{
	playerWidth := 5
	playerHeight := 5

	csWidth := 1
	csHeight := 1

	PlayerMinX := player.X - playerWidth
	PlayerMaxX := player.X + playerWidth
	PlayerMaxY := player.Y + playerHeight
	PlayerMinY := player.Y - playerHeight

	CSMinX := cs.X - csWidth
	CSMaxX := cs.X + csWidth
	CSMaxY := cs.Y + csHeight
	CSMinY := cs.Y - csHeight

	if PlayerMaxX < CSMinX|| PlayerMinX > CSMaxX {
		return false
	}

	if PlayerMaxY < CSMinY || PlayerMinY > CSMaxY {
		return false
	}
	return true
}

