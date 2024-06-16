package main

import (
	"fmt"
	"sync"
)

func main() {

	// UI goroutine
	uiTxRxChannel := make(chan string)
	uiCloseChannel := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		UiBootstrap(uiTxRxChannel, uiCloseChannel)
	}()

	// Keyboard goroutine
	go func() {
		BindHotkeys()
	}()

	// SDL2 goroutine

	uiTxRxChannel <- "Hello, UI thread!"

	fmt.Println("Main thread is doing other work")

	// Exit condition
	done := <-uiCloseChannel
	if done {
		fmt.Println("UI thread has finished")
	}

	// Wait on UI goroutine
	wg.Wait()
}
