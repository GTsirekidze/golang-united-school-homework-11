package batch

import (
	"fmt"
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
	c := make(chan user)
	var i int64 = 0
	var wg sync.WaitGroup
	go consume(c, &res)
	for ; i < pool; i++ {
		wg.Add(1)
		go getNUser(n, pool, i, c, &wg)
	}
	//fmt.Println("after wait")

	//fmt.Println("after close")

	wg.Wait()

	close(c)
	fmt.Println("HERE ends", res)
	return res
}

func consume(ch chan user, res *[]user) {
	for i := range ch {
		fmt.Println("HERE o", i)
		*res = append(*res, i)
		fmt.Println("HERE or not", res)
		//fmt.Println("HERE or not end", res)
	}
}

func getNUser(n int64, pool int64, routineIndex int64, userChan chan user, wg *sync.WaitGroup) {
	fmt.Println("123456789")
	var wtg sync.WaitGroup

	//for ; routineIndex <= n; routineIndex+=pool {
	for ; routineIndex < n; routineIndex += pool {
		wtg.Add(1)
		//fmt.Println("deadlock HERE or not", routineIndex)
		//go func(k int64, wtg *sync.WaitGroup) {
		//fmt.Println("ping routine K", k)
		userChan <- getOne(routineIndex)
		//fmt.Println("ping routine before  K done", k)
		wtg.Done()
		//fmt.Println("ping routine after  K done", k)
		//}(routineIndex, &wtg)
		//fmt.Println("deadlock free HERE or not", routineIndex)
	}
	wtg.Wait()
	wg.Done()
}
