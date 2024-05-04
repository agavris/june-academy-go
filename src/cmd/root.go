package cmd

import (
	"fmt"
	"github.com/agavris/june-academy-go/src/algorithm/scheduler"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Println(name, "took", time.Since(start))
	}
}

var rootCmd = &cobra.Command{
	Use:   "Scheduling",
	Short: "Run the scheduling algorithm.",
	Long:  `This is the entry point for the scheduling algorithm with a specified number of iterations.`,
	Run: func(cmd *cobra.Command, args []string) {
		Scheduler := scheduler.NewScheduler()
		defer timer("scheduling")()
		Scheduler.Run(numIterations)
	},
}

var numIterations int

func init() {
	rootCmd.PersistentFlags().IntVarP(&numIterations, "iterations", "n", 100, "Number of iterations to run the algorithm.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
