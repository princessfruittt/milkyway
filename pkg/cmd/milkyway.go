package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"milkyway/pkg/common/loggers"
)

// The Response value from Execute.
type Response struct {
	//// The build Result will only be set in the hugo build command.
	Result string

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
	commonCmd := newCommandBuilder().addAllCMDS().build()
	cmd := commonCmd.getCommand()
	cmd.SetArgs(args)

	c, err := cmd.ExecuteC()
	var response Response

	if c == cmd && commonCmd.c != nil {
		response.Result = "result"
	}
	if err == nil {
		errCount := int(loggers.GlobalErrorCounter.Count())
		if errCount > 0 {
			err = fmt.Errorf("logged %d errors", errCount)
		}
		//} else if response.Result != nil {
		//	errCount = response.Result.NumLogErrors()
		//	if errCount > 0 {
		//		err = fmt.Errorf("logged %d errors", errCount)
		//	}
		//}

	}

	response.Err = err
	response.Cmd = c

	return response
}

type result struct {
	res string
}

func (c *result) build() string {

	return "it is me"
}
