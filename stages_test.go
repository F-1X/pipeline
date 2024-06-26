package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStage1(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := -5; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	out := stage1(ch)
	b := []int{}
	for v := range out {
		b = append(b, v)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, b)
}

func TestStage2(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := -5; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	out := stage2(ch)
	var b []int
	for v := range out {
		b = append(b, v)
	}
	assert.Equal(t, []int{-5, -4, -2, -1, 1, 2, 4, 5, 7, 8}, b)
}

func TestStage3(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := -5; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	out := stage3(ch)
	var b []int
	for v := range out {
		b = append(b, v)
	}
	assert.Equal(t, []int{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, b)

}

func TestStage123(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := -5; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	out := stage3(stage2(stage1(ch)))
	var b []int
	for v := range out {
		b = append(b, v)
	}
	assert.Equal(t, []int{1, 2, 4, 5, 7, 8}, b)

}

func TestStage3Mod(t *testing.T) {
	stage3mod := Stage3Mod(stage3, 3, time.Second)

	ch := make(chan int)
	maxIter := 3
	go func() {
		for i := 0; i < maxIter; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		close(ch)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	var b []int
	go func() {
		go func() {
			time.Sleep(5 * time.Second)
			fmt.Println("Timeout: Closing Goroutine")
			wg.Done()
			return
		}()
		fmt.Println("start stage3mod")
		out := stage3mod(ch)

		for v := range out {
			fmt.Println("get v--->", v)
			b = append(b, v)
		}

	}()
	wg.Wait()

	assert.Equal(t, []int{0, 1, 2}, b)
}
