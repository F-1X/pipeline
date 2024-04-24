package main

type pipeline struct {
	stages stages
}

func NewPipeline(stages stages) *pipeline {
	return &pipeline{stages: stages}
}

func (p *pipeline) Run(ch chan int) chan int {

	out := make(chan int)

	for i := range p.stages {
		ch = p.stages[i](ch)
	}

	go func() {
		defer close(out)
		for v := range ch {
			out <- v
		}
	}()

	return out
}
