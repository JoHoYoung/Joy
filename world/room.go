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

var Rooms []Room
var NumberOfUser = 0

type Room struct {
	Players []string
	ClientMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
	ChanStart chan int
	Broadcast chan *Message
	Running bool
	Id int
}

func newRoom(id int) *Room{
	r := Room{}
	r.ClientMap = make(map[*Client]bool)
	r.ChanEnter = make(chan *Client)
	r.ChanLeave = make(chan *Client)
	r.Broadcast = make(chan *Message, conf.MESSAGE_BUFFER_SIZE * conf.USER_PER_ROOM)
	r.ChanStart = make(chan int,10)
	r.Id = id
	return &r
}

func (r *Room) broadCast(m *Message){
	fmt.Println("GET")
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
		fmt.Println("DELETE USER")
		delete(r.ClientMap, c)
		c.Delete()
	}
}

func (r *Room) run(){
	for {
		select{
		case  <-r.ChanEnter:
			if len(r.ClientMap) == conf.USER_PER_ROOM{
				fmt.Println("START")
				r.ChanStart <- 1
			}
			NumberOfUser++
		case c := <-r.ChanLeave:
			fmt.Println("LEAVE CLIENT")
			if _, ok := r.ClientMap[c]; ok {
				r.delUser(c)
			}
		case Message := <-r.Broadcast:
			r.broadCast(Message)
		case <- r.ChanStart:
			r.Running = true
			fmt.Println(len(r.ClientMap))
			players := make([]string, 0)
			for client, _ := range r.ClientMap{
				players = append(players, client.Id)
			}
			r.broadCast(&Message{Type:"START",Players:players})
			ChanTTL := time.NewTimer(time.Second * time.Duration(conf.PLAY_TIME_SEC))
			go func(){
				<-ChanTTL.C
				r.broadCast(&Message{Type:"END"})
				r.gameEnd()
			}()
		}
	}
}

func  (r *Room) Init(){
	for c := range r.ClientMap{
		c.Delete()
		delete(r.ClientMap, c)
	}
	r.Players = nil
	r.Running = false
}

func GenWord(){
	i:=0
	for i < conf.WORLDNUM {
		r := newRoom(i)
		Rooms = append(Rooms, *r)
		i++
		go r.run()
	}
}