package main

import (
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

var targetThrottle float64 = 0.0
var mu sync.Mutex

const deadzone = 5.0 // Define the deadzone threshold

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
					//fmt.Println("Throttle value received on controller:", msg)
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

		// Calculate the difference between the current value and the target value
		diff := currentValue - targetValue

		// Check if we should ignore the deadzone
		ignoreDeadzone := targetValue <= deadzone || targetValue >= 100-deadzone

		// Adjust throttle only if the difference is outside the deadzone,
		// or if we are near the extreme ends (0 or 100) and should ignore the deadzone
		if diff > deadzone && !ignoreDeadzone {
			robotgo.KeyDown("shift")
			robotgo.KeyUp("ctrl")
		} else if diff < -deadzone && !ignoreDeadzone {
			robotgo.KeyDown("ctrl")
			robotgo.KeyUp("shift")
		} else if ignoreDeadzone {
			if currentValue > targetValue {
				robotgo.KeyDown("shift")
				robotgo.KeyUp("ctrl")
			} else if currentValue < targetValue {
				robotgo.KeyDown("ctrl")
				robotgo.KeyUp("shift")
			}
		} else {
			// Within deadzone, ensure no keys are pressed
			robotgo.KeyUp("shift")
			robotgo.KeyUp("ctrl")
		}

		//fmt.Println("Current throttle:", currentValue, "Target throttle:", targetValue)

		// Wait a bit before checking again to avoid too rapid toggling
		time.Sleep(25 * time.Millisecond)
	}
}
