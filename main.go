package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dakshpareek/spine/cmd"
	"github.com/dakshpareek/spine/internal/types"
)

func main() {
	root := cmd.NewRootCmd()

	if err := root.Execute(); err != nil {
		exitCode := types.ExitCodeUserError

		var typed *types.Error

		if errors.As(err, &typed) && typed != nil {
			if typed.Code != 0 {
				exitCode = typed.Code
			}
			message := typed.Error()
			if message == "" && typed.Unwrap() != nil {
				message = typed.Unwrap().Error()
			}
			if message != "" {
				fmt.Fprintln(os.Stderr, message)
			}

		} else {
			fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(int(exitCode))
	}
}
