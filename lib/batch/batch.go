package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	//fmt.Println("getOne", id)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	c := make(chan user, n)
	var i int64 = 0
	var wg sync.WaitGroup
	for ; i < pool; i++ {
		wg.Add(1)
		go getNUser(n, pool, i, c, &wg)
	}
	wg.Wait()
	//fmt.Println("after wait")

	close(c)
	//fmt.Println("after close")

	for i := range c {
		//fmt.Println("HERE or not", i)
		//fmt.Println("HERE o", i)
		res = append(
			res, i)
		//fmt.Println("HERE or not end", res)
	}

	//fmt.Println("HERE negga", res)
	return res
}

func getNUser(n int64, pool int64, routineIndex int64, userChan chan user, wg *sync.WaitGroup) {
	//fmt.Println("123456789")
	var wtg sync.WaitGroup

	//for ; routineIndex <= n; routineIndex+=pool {
	for ; routineIndex < n; routineIndex += pool {
		wtg.Add(1)
		//fmt.Println("deadlock HERE or not", routineIndex)
		go func(k int64, wtg *sync.WaitGroup) {
			//fmt.Println("ping routine K", k)
			userChan <- getOne(k)
			//fmt.Println("ping routine before  K done", k)
			defer wtg.Done()
			//fmt.Println("ping routine after  K done", k)
		}(routineIndex, &wtg)
		//fmt.Println("deadlock free HERE or not", routineIndex)
	}
	wtg.Wait()
	wg.Done()
}
