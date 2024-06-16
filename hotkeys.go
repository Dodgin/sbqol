package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func BindHotkeys(uiTxRxChannel chan string) {

	oneshot := false

	hook.Register(hook.KeyDown, []string{"`", "alt"}, func(e hook.Event) {
		if !oneshot {
			fmt.Println("alt+` pressed")
			throttleInit()
			oneshot = true
		}

		// print throttle mappings
		fmt.Println("Throttle mappings:")
		for _, mapping := range ThrottleMappings {
			fmt.Println(mapping)
		}
		fmt.Println("End of throttle mappings")

		startThrottleWatcher()
	})

	s := hook.Start()
	<-hook.Process(s)
}
