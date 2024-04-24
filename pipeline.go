package main

import (
	"fmt"
)

type pipeline struct {
	stages stages
	// rb     *RingBuffer
}

func NewPipeline(stages stages) *pipeline {
	return &pipeline{stages: stages}
}

func (p *pipeline) Run(ch chan int) chan int {

	out := make(chan int)

	for i := range p.stages {
		ch = p.stages[i](ch)
	}
	fmt.Println("CH1",ch)
	go func() {
		defer close(out)
		fmt.Println("i liste")
		for v := range ch {
			fmt.Println("Get value in cha", v)
			out <- v
		}
	}()

	return out
}
