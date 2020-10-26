package faucet

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func cmdexec(bin string, args []string, inputs ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(bin, args...)
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
		return "", fmt.Errorf("error executing command: %s", stderr.String())
	}

	return strings.TrimSpace(stdout.String()), err
}

func (f *Faucet) cliexec(args []string, inputs ...string) (string, error) {
	return cmdexec(f.appCli, args, inputs...)
}
