package World

import (
	"fmt"
)
// 월드가 클라이언트 맵을
// 각 클라이언트가 자기 월드를.
// 월드가 클라이언트를 갖고있자.
// 클라이언트가 메세지를 받으면.
// | World |
// 1. 클라이언트가 메세지 받음.
// 2. 클라이언트가 자신이 속한 월드 채널에 메세지 날림.
// 3. 메세지를 받은 월드가 자신에 속한 클라이언트들에게 브로드 캐스팅.

var Rooms []Room
var NumberOfUser = 0

type Room struct {
	ClientMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
	Broadcast chan *Message
	Running bool
	Id int
}

func newRoom(id int) *Room{
	r := Room{}
	r.ClientMap = make(map[*Client]bool)
	r.ChanEnter = make(chan *Client)
	r.ChanLeave = make(chan *Client)
	r.Broadcast = make(chan *Message, config.MESSAGE_BUFFER_SIZE * config.USER_PER_ROOM)
	r.Id = id
	return &r
}

func (w *Room) broadCast(m *Message){
	fmt.Println("GET")
	for client, exist := range w.ClientMap{
		if(exist){
			client.Send <- m
		}
	}
}

func (r *Room) run(){
	for{
		select{
		case <-r.ChanEnter:
			fmt.Println("IN")
			NumberOfUser++;
		case c := <-r.ChanLeave:
			NumberOfUser--
			if _, ok := r.ClientMap[c]; ok {
				c.Delete()
			}
		case Message := <-r.Broadcast:
			r.broadCast(Message)
		}
	}
}

func  (r *Room) Init(){
	for c := range r.ClientMap{
		c.Delete()
	}
	r.Running = false;
}

func GenWord(){
	i:=0
	fmt.Println(config)
	for i < config.WORLDNUM {
		r := newRoom(i)
		Rooms = append(Rooms, *r)
		i++
		go r.run()
	}
}