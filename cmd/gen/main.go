package main

import (
	"fmt"
	"github.com/finderseyes/piper/cmd/gen/cmd"
	"os"
)

func main() {
	command := cmd.NewRootCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
