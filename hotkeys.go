package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func BindHotkeys(throttleChannel chan float64, uiTxRxChannel chan string) {

	oneshot := false

	hook.Register(hook.KeyDown, []string{"`", "alt"}, func(e hook.Event) {
		if !oneshot {
			fmt.Println("alt+` pressed")
			throttleInit()
			oneshot = true

			// print throttle mappings
			fmt.Println("Throttle mappings:")
			for _, mapping := range ThrottleMappings {
				fmt.Println(mapping)
				// print address in hex
				fmt.Printf("Address: %x\n", mapping.Address)
			}
			fmt.Println("End of throttle mappings")

			startThrottleWatcher(throttleChannel, uiTxRxChannel)
			ThrottleControllerStartMatching()
		}
	})

	s := hook.Start()
	<-hook.Process(s)
}
