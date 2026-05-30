package main

import (
	"Butter/inter"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a .but file to run.")
		return
	}

	inter.EachButterFunction()
	for name, body := range inter.ButterFunctions {
		if fn, ok := inter.NameToFunctions[name]; ok && fn.DoRun {
			inter.ButterInterpreter(body, inter.ButterLines[name], *fn, inter.ObjectFunc)
		}
	}
}
