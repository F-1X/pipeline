package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var sizeRingBuffer int = 3
var durationRingBuffer time.Duration = time.Second * 3

func main() {
	pipeline := NewPipeline(StagesBuilder(stage1, stage2, Stage3Mod(stage3, sizeRingBuffer, durationRingBuffer)))

	out := pipeline.Run(ReadNumbers())
	for v := range out {
		fmt.Println("v:", v)
	}

}

func ReadNumbers() chan int {
	ch := make(chan int)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "quit" || text == "exit" {
				close(ch)
				break
			}
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Error parsing input:", err)
				continue
			}
			ch <- num
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}

	}()
	return ch
}

// func WriteRandomNumbers()
