package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func readloop() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println(scanner.Text())

	}
}

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
	// заглушка, фейковые данные, ввод в канал
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		// input := []int{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5} // остануться 1, 2, 4, 5
		// for _, v := range input {
		// 	ch <- v
		// }

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
