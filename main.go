package main

import (
	"fmt"
	"sync"
)

func main() {

	// UI goroutine
	uiTxRxChannel := make(chan string)
	uiCloseChannel := make(chan bool)
	throttleChannel := make(chan float64, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		UiBootstrap(uiTxRxChannel, uiCloseChannel)
	}()

	// Keyboard goroutine
	go BindHotkeys(throttleChannel, uiTxRxChannel)

	// SDL2 goroutines
	go SdlBootstrap()
	go ThrottleControllerBootstrap(throttleChannel)

	// Debug message for testing UI channel
	uiTxRxChannel <- "{\"type\":\"debug\",\"value\":\"Hello, UI thread!\"}"

	// Debug message for testing main thread
	fmt.Println("Main thread is doing other work")

	// Exit condition
	done := <-uiCloseChannel
	if done {
		fmt.Println("UI thread has finished")
	}

	// Wait on UI goroutine
	wg.Wait()
}
