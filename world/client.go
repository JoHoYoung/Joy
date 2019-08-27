package world

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"joy/config"
	"sync"
)

// | Client |
// 1. 요청이오면 클라이언트를 만듬
// 2. world를 할당해서 월드에넣고, 클라이언트의 월드에 해당월드를 넣음.
// 3. 메세지가오면 월드로 보냄

var conf = config.Get()

var Mutex = &sync.Mutex{}
var pivot = 0

type Client struct {
	Conn *websocket.Conn // 웹소켓 커넥션
	Send chan *Message   // 메시지 전송용 채널
	Room *Room
	Id   string
}

func NewClient(conn *websocket.Conn) {
	uid, _ := uuid.NewUUID()
	c := &Client{
		Conn: conn,
		Send: make(chan *Message, conf.MESSAGE_BUFFER_SIZE),
		Id:   uid.String(),
	}
	c.AllocateWorld()
	go c.ReadLoop()
	go c.WriteLoop()
}

func (c *Client) Write(m *Message) {
	if _, ok := c.Room.ClientMap[c]; ok {
		c.Conn.WriteJSON(m)
	}
}

func (c *Client) WriteLoop() {
	for msg := range c.Send {
		c.Write(msg)
	}
}

func (c *Client) AllocateWorld() {
	Mutex.Lock()
	count := 0
	for (len(Rooms[pivot].ClientMap) >= conf.USER_PER_ROOM || Rooms[pivot].Running == true) {
		count ++
		pivot ++
		if pivot == len(Rooms) {
			pivot = 0
		}
		if (count == len(Rooms)) {
			c.Conn.Close()
			break
		}
	}
	c.Room = Rooms[pivot]
	c.Room.ChanEnter <- c
	Mutex.Unlock()
}

func (c *Client) Delete() {
	NumberOfUser--
	c.Conn.Close()
	fmt.Println("DELETE user, Id :", c.Id)
	close(c.Send)
}

func (c *Client) Read() (*Message, error) {
	var msg *Message
	if err := c.Conn.ReadJSON(&msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) ReadLoop() {
	for {
		m, err := c.Read()
		if err != nil { // 연결이 끊긴 클라이언트는 배열에서 제거..
			fmt.Println(err)
			c.Room.ChanLeave <- c
			break
		}
		c.Room.processMessage(*m)
	}
}
