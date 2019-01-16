package parallel

import (
	"fmt"
)

func main() {
	//var wg sync.WaitGroup
	var ch = make(chan struct{}, 10)

	seconds := [...]int{1,2,3,4,5}
	for i, s := range seconds {
		// 计数加 1
		//wg.Add(1)
		go func(i, s int) {
			// 计数减 1
			//defer wg.Done()
			fmt.Printf("goroutine%d 结束\n", i)
			ch <- struct{}{}
		}(i, s)
	}

	// 等待执行结束
	//wg.Wait()
	for i := 0; i < 5; i++ {
		<- ch
	}
	fmt.Println("所有 goroutine 执行结束")
}

//WaitGroup 用于等待一组 goroutine 结束ßß
//wg.Add() 方法一定要在 goroutine 开始前执行
//其实就是简化版的利用channel阻塞等待"子进程"跑完