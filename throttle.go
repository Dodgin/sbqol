package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type ThrottleAxes int

const (
	ThrottleAxis0 ThrottleAxes = iota
	ThrottleAxis1
	ThrottleAxis2
)

type ThrottleMapping struct {
	Name         string
	ThrottleAxis ThrottleAxes
	Address      uintptr
	Value        float32
}

var ThrottleMappings = []ThrottleMapping{}

func focusWindow(name string) bool {
	pids, err := robotgo.FindIds(name)
	if err != nil {
		fmt.Println("Error finding window IDs:", err)
		return false
	}

	if len(pids) > 0 {
		err = robotgo.ActivePid(pids[0])
		if err != nil {
			fmt.Println("Error focusing window:", err)
			return false
		}
		return true
	}

	fmt.Println("No window found with the name:", name)
	return false
}

func throttleInit() {
	if focusWindow("starbase.exe") {
		// Send w key down
		robotgo.Click()
		robotgo.KeyToggle("w", "down")

		initialstate, _ := getScanResults()

		ThrottleMappings = []ThrottleMapping{}

		// discard all initial values that are <= 0
		for _, value := range initialstate {
			if value.FloatValue > 0 {
				fmt.Println(value.FloatValue)
				ThrottleMappings = append(ThrottleMappings, ThrottleMapping{
					Name:         "FcuForward",
					ThrottleAxis: ThrottleAxis0,
					Address:      value.Address,
					Value:        value.FloatValue,
				})
			}
		}

		// print all throttle mappings
		for _, mapping := range ThrottleMappings {
			fmt.Println(mapping)
		}

		robotgo.KeyToggle("w", "up")
	} else {
		fmt.Println("Failed to focus the window.")
	}
}
