package main

import (
	"Butter/inter"
	"sync"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	var wg sync.WaitGroup

	inter.EachButterFunction()
	for key, line := range inter.ButterFunctions {
		if inter.NameToFunctions[key].DoRun {
			wg.Add(1)
			func() {
				inter.ButterInterpreter(line, inter.ButterLines[key], *inter.NameToFunctions[key], inter.ObjectFunc)
				wg.Done()
			}()
		}

	}
	wg.Wait()

}
