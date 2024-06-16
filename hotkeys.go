package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func BindHotkeys() {

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

		go func() {
			for {
				robotgo.MilliSleep(250)
				if len(ThrottleMappings) > 0 {
					GetFloat32ValueAtAddress(ThrottleMappings[0].Address)
					//fmt.Println("Fcuforward: ", val)
				}
			}
		}()
	})

	s := hook.Start()
	<-hook.Process(s)
}
