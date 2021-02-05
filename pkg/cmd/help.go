package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"regexp"
)

type cmdErrorHandler interface {
	getCommand() *cobra.Command
}

type cmdError struct {
	s         string
	userError bool
}

func (c cmdError) Error() string {
	return c.s
}

func (c cmdError) isUserError() bool {
	return c.userError
}
func newUserError(a ...interface{}) cmdError {
	return cmdError{s: fmt.Sprintln(a...), userError: true}
}

func newSystemError(a ...interface{}) cmdError {
	return cmdError{s: fmt.Sprintln(a...), userError: false}
}
func newSystemErrorFormatted(format string, a ...interface{}) cmdError {
	return cmdError{s: fmt.Sprintf(format, a...), userError: false}
}

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if error, ok := err.(cmdError); ok && error.isUserError() {
		return true
	}
	return userErrorRegexp.MatchString(err.Error())
}
