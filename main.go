package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now() //500ms
	task1 := task1()
	responseChannel := make(chan any, 2) //initialize channel of size 2
	
	waitGroup := &sync.WaitGroup{} // wait group that will wait for all tasks to finish
	waitGroup.Add(2)

	go task2(task1, responseChannel, waitGroup) //task2 will depend on task1
	go task3(task1, responseChannel,  waitGroup) //task3 will depend on task1

	waitGroup.Wait() //wait for all tasks to finish

	close(responseChannel) //close channel

	fmt.Println("Duration: ", time.Since(start))

	// Due to concurrency, task2 and task3 can be executed at the same time,
	// the total time of execution is less than the sum of the times all tasks,
	//T1: 500ms
	//T2: 550ms
	//T3: 600ms

	//T1: 500ms + T2: 550ms + T3: 600ms = 1.65s //without concurrency
	//T1: 500ms + max(T2: 550ms , T3: 600ms) = 1.1s //with concurrency
}

func task1() int {
	time.Sleep(500 * time.Millisecond) //lasts 500ms
	return 1
}
func task2( dependency int, resChan chan any, wg *sync.WaitGroup) {
	time.Sleep(550 * time.Millisecond) // lasts 550ms
	
	fmt.Println("T2 recieved dependency: ", dependency)

	resChan <- 1 //pass a return value to channel
	wg.Done() //inform wait group that task2 is done
}
func task3( dependency int, resChan chan any, wg *sync.WaitGroup) {
	time.Sleep(600 * time.Millisecond) // lasts 600ms

	fmt.Println("T3 recieved dependency: ", dependency)

	resChan <- 1
	wg.Done()
}
