package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {

	pipeline := NewPipeline(StagesBuilder(stage1, stage2, stage3))

	input := []int{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5} // остануться 1, 2, 4, 5
	expected := []int{1, 2, 4, 5}

	ch := make(chan int)
	go func() {
		for _, v := range input {
			ch <- v
		}
		close(ch)
	}()

	outCh := pipeline.Run(ch)

	actual := []int{}
	for v := range outCh {
		actual = append(actual, v)
	}

	assert.Equal(t, expected, actual)
}
