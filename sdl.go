package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var throttleValue float64 = 0.0

func GetThrottleValue() float64 {
	return throttleValue
}

func SdlBootstrap() {
	if err := sdl.Init(uint32(sdl.INIT_JOYSTICK)); err != nil {
		fmt.Printf("Failed to initialize SDL: %v\n", err)
		return
	}
	defer sdl.Quit()

	// List all connected joysticks
	numJoysticks := sdl.NumJoysticks()
	if numJoysticks < 1 {
		fmt.Println("No joysticks connected")
		return
	}
	fmt.Printf("Number of joysticks connected: %d\n", numJoysticks)
	for i := 0; i < numJoysticks; i++ {
		fmt.Printf("Joystick %d: %s\n", i, sdl.JoystickNameForIndex(i))
	}

	// Open the first joystick
	joystick := sdl.JoystickOpen(1)
	if joystick == nil {
		fmt.Printf("Could not open joystick 0: %v\n", sdl.GetError())
		return
	}
	defer joystick.Close()

	// Goroutine to read joystick axis
	go readJoystickAxis(joystick)

	fmt.Println("SDL2 initialized")

	// Keep the function running
	select {}
}

func convertToFloatScale(value int16) float64 {
	// Ensure the value is within the expected range
	if value < -32768 || value > 32767 {
		panic("value out of range")
	}

	// Convert the range [-32768, 32767] to [0, 1]
	// 32767 should map to 0
	// -32768 should map to 1
	floatValue := 1.0 - (float64(value)+32768.0)/65535.0

	return floatValue
}

func readJoystickAxis(joystick *sdl.Joystick) {
	fmt.Println("Reading joystick axis")
	for {
		sdl.PumpEvents() // Update the state of all connected devices

		// Directly read the joystick axis
		axisValue := joystick.Axis(2)
		//fmt.Printf("Joystick 0 Axis 0 value: %f\n", convertToFloatScale(axisValue))
		throttleValue = convertToFloatScale(axisValue)

		time.Sleep(25 * time.Millisecond) // Poll every 100ms
	}
}
