package main

import (
	"github.com/spf13/cobra"
	"AnsiToTosca/pkg/cmd/help.go"
)

// The Response value from Execute.
type Response struct {
	//// The build Result will only be set in the hugo build command.
	//Result *hugolib.HugoSites

	// Err is set when the command failed to execute.
	Err error

	// The command that was executed.
	Cmd *cobra.Command
}

// IsUserError returns true is the Response error is a user error rather than a
// system error.
func (r Response) IsUserError() bool {
	return r.Err != nil && isUserError(r.Err)
}

func Execute(args []string) Response {
	commonCmd
}
