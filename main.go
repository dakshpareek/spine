package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dakshpareek/ctx/cmd"
	"github.com/dakshpareek/ctx/internal/types"
)

var version = "dev"

func main() {
	root := cmd.NewRootCmd(version)

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
