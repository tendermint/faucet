package faucet

func (f *Faucet) loadKey() error {
	if !f.keyExists(f.keyName) && f.keyMnemonic != "" {
		if _, err := f.cliexec([]string{"keys", "add", f.keyName, "--recover"}, f.keyMnemonic, f.keyringPassword, f.keyringPassword); err != nil {
			return err
		}
	}

	var err error
	f.faucetAddress, err = f.cliexec([]string{"keys", "show", f.keyName, "-a"}, f.keyringPassword)
	if err != nil {
		return err
	}

	return nil
}

func (f *Faucet) keyExists(keyname string) bool {
	if _, err := f.cliexec([]string{"keys", "show", keyname}, f.keyringPassword); err == nil {
		return true
	}
	return false
}
