package parallel

import (
	"fmt"
	"sync"
	"time"
	"github.com/go-redis/redis"
)

func incr() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var lockKey = "counter_lock"
	var counterKey = "counter"

	// lock
	resp := client.SetNX(lockKey, 1, time.Second*5)
	lockSuccess, err := resp.Result()

	if err != nil || !lockSuccess {
		fmt.Println(err, "lock result: ", lockSuccess)
		return
	}

	// counter ++
	getResp := client.Get(counterKey)
	cntValue, err := getResp.Int64()
	if err == nil {
		cntValue++
		resp := client.Set(counterKey, cntValue, 0)
		_, err := resp.Result()
		if err != nil {
			// log err
			println("set value error!")
		}
	}
	println("current counter is ", cntValue)

	delResp := client.Del(lockKey)
	unlockSuccess, err := delResp.Result()
	if err == nil && unlockSuccess > 0 {
		println("unlock success!")
	} else {
		println("unlock failed", err)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incr()
		}()
	}
	wg.Wait()
}

//核心：分布式节点通过对redis的同一个key-value的维护实现一个lock的概念
//知识点：
//1. redis setNX支持设置keyvalue时指定expiretime，避免非原子操作时，刚加了锁线程挂了，锁永远被锁住
//			setNX成功时，返回1，否则0.（按照db操作时，产生影响的行数来记）这样保证其他线程setNX的时候能够根据resp判断是否被锁了
//			set value时，可以设置线程id，删之前加一个校验，避免加锁线程没执行完超时了，另外一个加了锁，等第一个加锁的线程去释放的时候释放的是第二个线程加的锁--存在问题：两个线程同时访问了共享数据
//2. 考虑守护进程：强行续命，保证是线程主动del而不是超时释放

//其他分布式锁，找时间研究：
//zk、memcached， etcd， redlock， 分析优缺点，关注使用场景