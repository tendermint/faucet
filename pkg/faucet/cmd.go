package faucet

import (
	"bytes"
	"fmt"
	"os/exec"
)

func (f *Faucet) executeCli(args []string, inputs ...string) (string, error) {
	cmd := exec.Command(f.appCli, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	for _, input := range inputs {
		if _, err := fmt.Fprintln(stdin, input); err != nil {
			return "", err
		}
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("error executing command: %s", string(stderr.Bytes()))
	}

	return string(stdout.Bytes()), err
}
