package main

import (
	"fmt"
	"time"
)

type stage func(chan int) chan int

type stages []stage

func StagesBuilder(x ...stage) stages {
	stages := make([]stage, 0)
	stages = append(stages, x...)
	return stages
}

func stage1(ch chan int) chan int {
	chOut := make(chan int)
	go func() {
		for i := range ch {
			if i >= 0 {
				chOut <- i
			}
		}
		close(chOut)
	}()
	return chOut
}

func stage2(ch chan int) chan int {
	chOut := make(chan int)
	go func() {
		for i := range ch {
			if i%3 != 0 && i != 0 {
				chOut <- i
			}
		}
		close(chOut)
	}()

	return chOut
}

func stage3(ch chan int) chan int {
	outCh := make(chan int)
	go func() {
		for i := range ch {
			// fmt.Println("stage3 get i",i)
			outCh <- i
		}
		close(outCh)
	}()
	return outCh
}

// обертка для stage3 дополняющая функционалом RingBuffer, в промежутчной коммуниции каналов
func Stage3Mod(fun func(ch chan int) chan int, size int, seconds time.Duration) func(ch chan int) chan int {
	rb := NewRingBuffer(size, seconds)
	go rb.loop()

	return func(ch chan int) chan int {
		
		go func() {
			// defer close(rb.tmp) // некоторые сложности с синхронизацией отпуска по времени но канал по идее закрывать надо
			for i := range stage3(ch) {
				rb.putCh <- i
			}
			fmt.Println("Close putch")
			rb.done <- struct{}{}
			close(rb.putCh)
		}()
		return rb.tmp
	}

}
