package faucet

import "strings"

func (f *Faucet) loadKey() error {
	if !f.keyExists(f.keyName) && f.keyMnemonic != "" {
		if _, err := f.executeCli([]string{"keys", "add", f.keyName, "--recover"}, f.keyMnemonic, f.keyringPassword, f.keyringPassword); err != nil {
			return err
		}
	}

	output, err := f.executeCli([]string{"keys", "show", f.keyName, "-a"}, f.keyringPassword)
	if err != nil {
		return err
	}
	f.faucetAddress = strings.TrimSpace(output)

	return nil
}

func (f *Faucet) keyExists(keyname string) bool {
	if _, err := f.executeCli([]string{"keys", "show", keyname}, f.keyringPassword); err == nil {
		return true
	}
	return false
}
