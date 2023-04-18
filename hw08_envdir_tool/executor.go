package main

import (
	"os"
	"os/exec"
)

const exitError = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return exitError
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := handleEnvVars(env); err != nil {
		return exitError
	}

	_ = command.Run()

	return command.ProcessState.ExitCode()
}

func handleEnvVars(env Environment) error {
	for envVar, val := range env {
		if val.NeedRemove {
			if err := os.Unsetenv(envVar); err != nil {
				return err
			}
			break
		}

		if _, ok := os.LookupEnv(envVar); ok {
			if err := os.Unsetenv(envVar); err != nil {
				return err
			}
		}

		if err := os.Setenv(envVar, val.Value); err != nil {
			return err
		}
	}

	return nil
}
