package main

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm/scheduler"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Println(name, "took", time.Since(start))
	}
}

func main() {

	Scheduler := scheduler.NewScheduler()
	defer timer("main")()
	Scheduler.Run(10000)

}
