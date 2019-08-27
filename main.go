package main

// 8명 채우고, 다음 월드, 채우고 다음 월드.
// 동시에들어오면... Critical Section문제.. 한 월드에 100명 몰릴수도 있음
// mutex걸어야하나? 걸면 느려질것 같은데
// 안걸고 할수있는 방법은 없을까??...
// 월드는 공유되고, 각 월드에 몇명은 공유자원임.
// 한명한명 올때마다 락걸고  +1 할당하고 .. 느려질거같은데

import (
	"joy/v1"
	"joy/world"
)

func main() {
	world.GenWord()
	v1.Start()
}
