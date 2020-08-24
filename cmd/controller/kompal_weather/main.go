package main

import (
	"fmt"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s\n", err)
			os.Exit(1)
		}
	}()
	//cli := &gig.Gig{OutStream: os.Stdout, ErrStream: os.Stderr, Version: version}
	//os.Exit(cli.Run())
}
