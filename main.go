package main

import (
	"fmt"
	"os"

	"github.com/Pixie2468/git-branch-explorer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
