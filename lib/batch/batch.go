package batch

import (
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
	var wg, wg2 sync.WaitGroup
	wg.Add(int(pool))
	wg2.Add(1)
	go consume(c, &res, &wg2)

	for ; i < pool; i++ {
		go getNUser(n, pool, i, c, &wg)
	}

	wg.Wait()
	close(c)
	wg2.Wait()
	return res
}

func consume(ch chan user, res *[]user, wg *sync.WaitGroup) {

	for i := range ch {
		Mu.Lock()
		*res = append(*res, i)
		Mu.Unlock()
	}
	defer wg.Done()

}

func getNUser(n int64, pool int64, routineIndex int64, userChan chan user, wg *sync.WaitGroup) {
	k := routineIndex
	for ; k < n; k += pool {
		userChan <- getOne(k)
	}
	wg.Done()
}
