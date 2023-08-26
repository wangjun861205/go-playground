package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	group := singleflight.Group{}
	serial(&group)
	parallel_without_overlap(&group)
	parallel_with_overlap(&group)
}

func serial(group *singleflight.Group) {
	res1, _, isMemo1 := group.Do("test", func() (interface{}, error) {
		return "result 1", nil
	})
	fmt.Printf("result: %v, is_memo: %t\n", res1, isMemo1)
	res2, _, isMemo2 := group.Do("test", func() (interface{}, error) {
		return "result 2", nil
	})
	fmt.Printf("result: %v, is_memo: %t\n", res2, isMemo2)
}

func parallel_without_overlap(group *singleflight.Group) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res1, _, isMemo1 := group.Do("test", func() (interface{}, error) {
			return "result 1", nil
		})
		fmt.Printf("result: %v, is_memo: %t\n", res1, isMemo1)
	}()
	time.Sleep(500 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		res1, _, isMemo1 := group.Do("test", func() (interface{}, error) {
			return "result 2", nil
		})
		fmt.Printf("result: %v, is_memo: %t\n", res1, isMemo1)
	}()
	wg.Wait()

}

func parallel_with_overlap(group *singleflight.Group) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res1, _, isMemo1 := group.Do("test", func() (interface{}, error) {
			time.Sleep(500 * time.Millisecond)
			return "result 1", nil
		})
		fmt.Printf("result: %v, is_memo: %t\n", res1, isMemo1)
	}()
	time.Sleep(100 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		res1, _, isMemo1 := group.Do("test", func() (interface{}, error) {
			return "result 1", nil
		})
		fmt.Printf("result: %v, is_memo: %t\n", res1, isMemo1)
	}()
	wg.Wait()

}
