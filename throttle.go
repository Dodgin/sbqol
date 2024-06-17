package main

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/go-vgo/robotgo"
)

type ThrottleAxes int

const (
	ThrottleAxis0 ThrottleAxes = iota
	ThrottleAxis1
	ThrottleAxis2
)

type ThrottleMapping struct {
	Name         string       `json:"name"`
	ThrottleAxis ThrottleAxes `json:"throttleAxis"`
	Address      uintptr      `json:"address"`
	Value        float32      `json:"value"`
}

var ThrottleMappings = []ThrottleMapping{}

func startThrottleWatcher(throttleChannel chan float64, uiTxRxChannel chan string) {
	go func() {
		for {
			robotgo.MilliSleep(50)
			if len(ThrottleMappings) > 0 {
				val, _ := GetFloat32ValueAtAddress(ThrottleMappings[0].Address)
				ThrottleMappings[0].Value = val
				jsonData, _ := json.Marshal(ThrottleMappings)
				jsonMessage := UiMessage{
					Type:    "throttle",
					Payload: string(jsonData),
				}
				jsonMessageData, _ := json.Marshal(jsonMessage)

				// Notify UI of new value
				uiTxRxChannel <- string(jsonMessageData)

				// Notify Throttle controller
				//fmt.Println("Throttle value received on watcher:", val)
				select {
				case throttleChannel <- float64(val):
					// Value sent successfully
				default:
					// Channel is full, read and discard the old value
					<-throttleChannel
					// Now send the new value
					throttleChannel <- float64(val)
				}

			}
		}
	}()
}

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
			//cast to int and chack if it is greater than 0
			if math.Floor(float64(value.FloatValue)) > 0 {
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
