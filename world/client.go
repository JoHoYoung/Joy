package world

import (
	"fmt"
	"github.com/gorilla/websocket"
	"joy/config"
	"log"
	"sync"
)

// | Client |
// 1. 요청이오면 클라이언트를 만듬
// 2. world를 할당해서 월드에넣고, 클라이언트의 월드에 해당월드를 넣음.
// 3. 메세지가오면 월드로 보냄

var conf = config.Get();
var Mutex = &sync.Mutex{}
var pivot = 0

var Clients []*Client

type Client struct {
	Conn *websocket.Conn // 웹소켓 커넥션
	Send chan *Message   // 메시지 전송용 채널
	Room *Room
	User   *User  // 현재 접속한 사용자 정보
}

func (c *Client) Write(m *Message) error {
	log.Println("write to websocket:", m)
	return c.Conn.WriteJSON(m)
}

func (c *Client) WriteLoop() {
	for msg := range c.Send {
		c.Write(msg)
	}
}


func (c *Client) AllocateWorld(){
	Mutex.Lock()
	count := 0
	for len(Rooms[pivot].ClientMap) >= conf.USER_PER_ROOM || Rooms[pivot].Running {
		// Out logic
		count ++
		pivot ++
		fmt.Println("count", count)
		if(count == len(Rooms)){
			c.Conn.Close()
			break;
		}
	}
	fmt.Println("IN",pivot)
	Rooms[pivot].ClientMap[c] = true
	c.Room = &Rooms[pivot];
	c.Room.ChanEnter <- c
	Mutex.Unlock()
}

func NewClient(conn *websocket.Conn, u *User){
	c := &Client{
		Conn: conn,
		Send: make(chan *Message, conf.MESSAGE_BUFFER_SIZE),
		User: u,
	}
	c.AllocateWorld()
	go c.ReadLoop()
	go c.WriteLoop()
}

func (c *Client) Delete(){
	delete(c.Room.ClientMap, c)
	close(c.Send)
}

func (c *Client) Read() (*Message, error) {
	var msg *Message
	if err := c.Conn.ReadJSON(&msg); err != nil {
		fmt.Println("ERRR")
		return nil, err
	}
	log.Println("read from websocket:", msg)
	return msg, nil
}

func (c *Client) ReadLoop() {

	defer func() {
		c.Room.ChanLeave <- c
		c.Conn.Close()
	}()

	for {
		m, err := c.Read()
		if err != nil { // 연결이 끊긴 클라이언트는 배열에서 제거..
			log.Println("read message error: ", err)
			break
		}
		fmt.Println("BROAD")
		c.Room.Broadcast <- m
	}
}