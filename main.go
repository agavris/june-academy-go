package main

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/cmd"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("cpu_prof_ja.prof")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		return
	}
	defer pprof.StopCPUProfile()

	cmd.Execute()
	fmt.Println("Executed...")
	//cmd.Execute()
}
