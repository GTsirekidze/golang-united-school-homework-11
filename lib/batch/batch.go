package batch

import (
	"fmt"
	"sync"
	"time"
)

type user struct {
	ID int64
}

var Mu sync.Mutex

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	c := make(chan user)
	var i int64 = 0
	var wg sync.WaitGroup
	wg.Add(int(pool))

	go consume(c, &res)

	for ; i < pool; i++ {

		go getNUser(n, pool, i, c, &wg)
	}

	wg.Wait()

	close(c)

	fmt.Println("HERE ends everything", res)
	return res
}

func consume(ch chan user, res *[]user) {

	//fmt.Println("HERE ends", res, <-ch)
	//fmt.Println("index is")
	for i := range ch {
		Mu.Lock()
		*res = append(*res, i)
		Mu.Unlock()
		//fmt.Println("HERE ends", res, i)
	}

}

func getNUser(n int64, pool int64, routineIndex int64, userChan chan user, wg *sync.WaitGroup) {
	//fmt.Println("123456789")
	k := routineIndex
	for ; k < n; k += pool {

		userChan <- getOne(k)

	}
	wg.Done()
}
