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
	command.Env = buildEnv(env)

	_ = command.Run()

	return command.ProcessState.ExitCode()
}

func buildEnv(env Environment) []string {
	for key, env := range env {
		if env.NeedRemove {
			_ = os.Unsetenv(key)
			continue
		}

		_ = os.Setenv(key, env.Value)
	}

	return os.Environ()
}
