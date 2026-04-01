package main

import (
	"go-upcycle_connect-backend/cmd"
	"os"
)

func main() {
	if len(os.Args) > 0 {
		cmd.Cmd()
	}
}
