package world

import (
	"fmt"
	"time"
)
// 월드가 클라이언트 맵을
// 각 클라이언트가 자기 월드를.
// 월드가 클라이언트를 갖고있자.
// 클라이언트가 메세지를 받으면.

// | world |
// 1. 클라이언트가 메세지 받음.
// 2. 클라이언트가 자신이 속한 월드 채널에 메세지 날림.
// 3. 메세지를 받은 월드가 자신에 속한 클라이언트들에게 브로드 캐스팅.

var Rooms []*Room
var NumberOfUser = 0


type Room struct {
	ClientMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
	ChanStart chan int
	Broadcast chan *Message
	Running bool
	Id int
	Players map[string]*Player
}

func newRoom(id int) *Room{
	r := Room{}
	r.Players = make(map[string]*Player)
	r.ClientMap = make(map[*Client]bool)
	r.ChanEnter = make(chan *Client,10)
	r.ChanLeave = make(chan *Client,10)
	r.Broadcast = make(chan *Message, conf.MESSAGE_BUFFER_SIZE * conf.USER_PER_ROOM)
	r.ChanStart = make(chan int,10)
	r.Running = false
	r.Id = id
	return &r
}

func (r *Room) broadCast(m *Message){
	for client, exist := range r.ClientMap{
		if(exist){
			client.Send <- m
		}
	}
}

func (r *Room) gameEnd(){
	r.Init()
}

func (r *Room) delUser(c *Client){
	if _, ok := r.ClientMap[c]; ok{
		delete(r.ClientMap, c)
		c.Delete()
	}
}

func (r *Room) regenEvent(){
	ticker := time.NewTicker(time.Second * time.Duration(conf.REGEN_TIME_SEC))
	for{
		<- ticker.C
		if !r.Running{
			break
		}
		r.Broadcast <- &Message{Type:conf.TYPE.REGEN}
	}
}

func (r *Room) dataEvent(){
	ticker := time.NewTicker(time.Second * time.Duration(conf.SEND_TIME_SEC))
	for{
		<- ticker.C
		r.simulate()
		if !r.Running{
			break
		}
		r.Broadcast <- &Message{Type:conf.TYPE.PLAY,Players: r.Players}
	}
}

func (r *Room) endEvent(){
	ChanTTL := time.NewTimer(time.Second * time.Duration(conf.PLAY_TIME_SEC))
	<-ChanTTL.C
	r.broadCast(&Message{Type:conf.TYPE.END})
	r.gameEnd()
}

func (r *Room) run(){
	for {
		select {
		case c := <-r.ChanEnter:
			r.ClientMap[c] = true
			r.Players[c.Id] = NewPlayer(c.Id)
			c.Send <- &Message{Type: conf.TYPE.ENTER, Player: c.Room.Players[c.Id]}
			NumberOfUser++
			if len(r.ClientMap) == conf.USER_PER_ROOM {
				r.ChanStart <- 1
			}
		case c := <-r.ChanLeave:
			fmt.Println("Client leave, ID : ", c.Id)
			if _, ok := r.ClientMap[c]; ok {
				r.delUser(c)
			}
		case Message := <-r.Broadcast:
			r.broadCast(Message)
		case <-r.ChanStart:
			r.Running = true
			r.broadCast(&Message{Type: conf.TYPE.START, Players: r.Players, Score:r.GetScore()})
			go r.dataEvent()
			go r.regenEvent()
			go r.endEvent()
		}
	}
}

func (r *Room)simulate(){
	fmt.Println("SIMULATE")
	for _, M :=  range r.Players{
		for _, S := range r.Players{
			if M.ID != S.ID{
				for i,cs := range S.CS{
					if collision(*M, cs){
						upper := len(M.CS)
						M.CS = append(M.CS, S.CS[i:]...)
						S.CS = S.CS[0:i]
						M.GetCS = len(M.CS) - upper
						break
					}
				}
			}
		}
	}
}

func (r *Room) Init(){
	for c := range r.ClientMap{
		fmt.Println("Room initialize, ID :",c.Id)
		c.Conn.Close()
	}
	r.Players = make(map[string]*Player)
	r.Running = false
}

func (r *Room) GetScore() map[string]int{
	result := make(map[string]int,10)
	for _, c := range r.Players {
		result[c.ID] = c.Score
	}
	return result
}

func (r *Room) processMessage(m Message){
	if m.Type == conf.TYPE.START{
		if _, ok := r.Players[m.Player.ID]; ok{
			r.Players[m.Player.ID] = m.Player
		}
	}
	if m.Type == conf.TYPE.IN{
		if _, ok := r.Players[m.Player.ID]; ok{
			r.Players[m.Player.ID].Score += m.Player.GetCS
			r.Players[m.Player.ID].GetCS = 0
		}
	}
}

func GenWord(){
	i:=0
	for i < conf.WORLDNUM {
		r := newRoom(i)
		Rooms = append(Rooms, r)
		i++
		go r.run()
	}
}