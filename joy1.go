package main

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/go-vgo/robotgo"
)

type JoyAxes int

const (
	JoyAxis0 JoyAxes = iota
	JoyAxis1
	JoyAxis2
)

type JoyMapping struct {
	Name    string  `json:"name"`
	JoyAxis JoyAxes `json:"joyAxis"`
	Address uintptr `json:"address"`
	Value   float32 `json:"value"`
}

var Joy1Mappings = []JoyMapping{}

func startJoyWatcher(joyChannel chan float64, uiTxRxChannel chan string) {
	go func() {
		for {
			robotgo.MilliSleep(50)
			if len(Joy1Mappings) > 0 {
				val, _ := GetFloat32ValueAtAddress(Joy1Mappings[0].Address)
				Joy1Mappings[0].Value = val
				jsonData, _ := json.Marshal(Joy1Mappings)
				jsonMessage := UiMessage{
					Type:    "joy",
					Payload: string(jsonData),
				}
				jsonMessageData, _ := json.Marshal(jsonMessage)

				// Notify UI of new value
				uiTxRxChannel <- string(jsonMessageData)

				// Notify Joy controller
				//fmt.Println("Joy value received on watcher:", val)
				select {
				case joyChannel <- float64(val):
					// Value sent successfully
				default:
					// Channel is full, read and discard the old value
					<-joyChannel
					// Now send the new value
					joyChannel <- float64(val)
				}

			}
		}
	}()
}

func joy1Init() {
	if focusWindow("starbase.exe") {

		Joy1Mappings = []JoyMapping{}

		robotgo.Click()

		// Map pitch
		robotgo.KeyToggle("s", "down")
		initialstate, _ := getScanResults()
		for _, value := range initialstate {
			if math.Floor(float64(value.FloatValue)) > 0 {
				fmt.Println(value.FloatValue)
				Joy1Mappings = append(Joy1Mappings, JoyMapping{
					Name:    "FcuRotationalPitch",
					JoyAxis: JoyAxis0,
					Address: value.Address,
					Value:   value.FloatValue,
				})
			}
		}
		robotgo.KeyToggle("s", "up")

		robotgo.MilliSleep(500)

		// Map roll
		robotgo.KeyToggle("d", "down")
		initialstate, _ = getScanResults()
		for _, value := range initialstate {
			if math.Floor(float64(value.FloatValue)) > 0 {
				fmt.Println(value.FloatValue)
				Joy1Mappings = append(Joy1Mappings, JoyMapping{
					Name:    "FcuRotationalRoll",
					JoyAxis: JoyAxis0,
					Address: value.Address,
					Value:   value.FloatValue,
				})
			}
		}
		robotgo.KeyToggle("d", "up")

		robotgo.MilliSleep(500)

		// Map yaw
		robotgo.KeyToggle("e", "down")
		initialstate, _ = getScanResults()
		for _, value := range initialstate {
			if math.Floor(float64(value.FloatValue)) > 0 {
				fmt.Println(value.FloatValue)
				Joy1Mappings = append(Joy1Mappings, JoyMapping{
					Name:    "FcuRotationalYaw",
					JoyAxis: JoyAxis0,
					Address: value.Address,
					Value:   value.FloatValue,
				})
			}
		}
		robotgo.KeyToggle("e", "up")

		// print all joy mappings
		for _, mapping := range Joy1Mappings {
			fmt.Println(mapping)
		}

	}
}
