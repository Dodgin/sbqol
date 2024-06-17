package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

var targetThrottle float64 = 0.0
var mu sync.Mutex

func ThrottleControllerStartMatching() {
	// Start the throttle adjustment goroutine
	go adjustThrottle()
}

// This sends throttle inputs to the game
func ThrottleControllerBootstrap(throttleChannel chan float64) {
	go func() {
		for {
			select {
			case msg := <-throttleChannel:
				if msg != targetThrottle {
					fmt.Println("Throttle value received on controller:", msg)
					updateTargetThrottle(msg)
				}
			}
		}
	}()
}

// Get the current target throttle level
func getTargetThrottle() float64 {
	mu.Lock()
	defer mu.Unlock()
	return targetThrottle
}

// Update the target throttle level
func updateTargetThrottle(value float64) {
	mu.Lock()
	defer mu.Unlock()
	targetThrottle = value
}

func adjustThrottle() {

	for {
		currentValue := GetThrottleValue() * 100
		targetValue := getTargetThrottle()

		if currentValue > targetValue {
			robotgo.KeyDown("w")
			robotgo.KeyUp("s")
		} else if currentValue < targetValue {
			robotgo.KeyDown("s")
			robotgo.KeyUp("w")
		} else {
		}

		fmt.Println("Current throttle:", currentValue, "Target throttle:", targetValue)

		// Wait a bit before checking again to avoid too rapid toggling
		time.Sleep(25 * time.Millisecond)
	}
}
