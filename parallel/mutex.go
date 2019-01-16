package parallel

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan struct{}, 2)

	var l sync.Mutex
	go func() {
		l.Lock()
		defer l.Unlock()
		fmt.Println("goroutine1: 我会锁定大概 2s")
		time.Sleep(time.Second * 2)
		fmt.Println("goroutine1: 我解锁了，你们去抢吧")
		ch <- struct{}{}
	}()

	go func() {
		fmt.Println("groutine2: 等待解锁")
		l.Lock()
		defer l.Unlock()
		fmt.Println("goroutine2: 哈哈，我锁定了")
		ch <- struct{}{}
	}()

	// 等待 goroutine 执行结束
	for i := 0; i < 2; i++ {
		<-ch
	}
}



//锁其实锁的是对应的锁，不同goroutine的对同一个锁进行lock时候，如果没办法锁住，那么就会等能锁住的时候继续往下走-block

//LINE32: 带缓存的channel是为了防止主goroutine提前退出后，两个起来的goroutine还没结束导致看不到效果，
// 同理可以sleep，同时receive两次必然两个goroutine要send两次，即跑完两个go

//双重lock 报错是因为lock是在同一个goroutine中串行执行的，没解锁必然不能lock

//Q：goroutine是线程还是进程
//Q:协程是什么——flask中的协程库及当时存在的问题