package api

import "time"

func Add(a, b int) int {
	return a + b
}

func Tick() chan int {
	c := make(chan int)
	go func() {
		for {
			time.Sleep(time.Second)
			c <- 1
		}
	}()
	return c
}
