package main

import (
	"github.com/princessfruittt/ansitotosca/pkg/cmd"
	"os"
)

func main() {
	resp := cmd.Execute(os.Args[1:])

	if resp.Err != nil {
		if resp.IsUserError() {
			resp.Cmd.Println("")
			resp.Cmd.Println(resp.Cmd.UsageString())
		}
		os.Exit(-1)
	}
}
