### JOY 
Websocket game
#### BackEnd : Go(gorilla + gin) | FrontEnd : Cocos2D 

### 클라이언트 ( Socket )접속
![Mutex](https://user-images.githubusercontent.com/37579650/68004819-485fdc00-fcb6-11e9-8fd6-b60fd774395b.png)

* 유저가 접속하면 현재 접속가능한 방을 탐색한다( 인원수가 다 차지않은방, 실행하고 있지 않은 방)
* 스레딩시, 유저인원수 제한을 넘을 수 있으므로 Mutex를 사용해 보호한다.
* 비어있는 방을 찾으면, 유저를 해당 방에 할당하고 Mutex를 반납한다
```
func (c *Client) AllocateWorld() {

	Mutex.Lock()
	defer Mutex.Unlock()

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
}

```

### 통신
![Untitled Diagram](https://user-images.githubusercontent.com/37579650/68004820-485fdc00-fcb6-11e9-84c8-b706a2712701.png)

* 유저가 방에 접속하면 모든 유저의 통신은 방을 통한다.
* 유저 객체는 자신이속한 방을 Referencing 한다.
```
func (c *Client) AllocateWorld() {
    
    .
    .
    .
	c.Room = Rooms[pivot]
	c.Room.ChanEnter <- c
}
```
* 힙에 유저객체를 할당하고, 읽기 쓰기 고루틴을 생성한다.
```
func NewClient(conn *websocket.Conn) {
	uid, _ := uuid.NewUUID()
	var c = &Client{
		Conn: conn,
		Send: make(chan *Message, conf.MESSAGE_BUFFER_SIZE),
		Id:   uid.String(),
	}
	c.AllocateWorld()
	go c.ReadLoop()
	go c.WriteLoop()
}
```
* 소켓으로 액션이 들어오면, 자신이 속한 방으로 emit한다.
```
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

```
* 방은 액션을 받으면 자신의 클라이언트 들에게 broadcasting 한다.
```
func (r *Room) broadCast(m *Message){
	for client, exist := range r.ClientMap{
		if(exist){
			client.Send <- m
		}
	}
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
```