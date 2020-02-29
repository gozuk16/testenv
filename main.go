package main

import (
	"fmt"
	"os"

	"github.com/gozuk16/gret/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
