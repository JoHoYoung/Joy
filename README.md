### JOY 
#### WebSocket io Games made with Golang and coke2D. Each user eats a neutral object on the map and attaches it to his tail, and when he enters his camp, he can score as many points as he eats. Users can intercept a tail that is worn by another user by hitting the tail of another
* * *
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

### 동기화
* 실시간 게임 이기 떄문에 동기화 문제를 해결해야함
* 각 클라이언트에서 cocos2D 물리엔진으로 충돌판정을 하면 동기화문제가 있을 수 있다고 판단
* 서버에서 AABB를 판단
```
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

```
* 60 fps를 기준으로 서버에서 AABB로 충돌을 판정, 특정 타임틱마다 같은 충돌 판정 결과 즉 같은 데이터를 동시에 뿌려줌
