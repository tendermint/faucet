package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/allinbits/cosmos-faucet/pkg/faucet"
	"github.com/allinbits/cosmos-faucet/pkg/utils"
)

var (
	logLevel string
	port     int

	keyName         string
	keyMnemonic     string
	keyringPassword string
	appCli          string
	denom           string
	creditAmount    uint64
	maxCredit       uint64
)

func init() {
	flag.StringVar(&logLevel, "log-level", utils.GetEnvString("LOG_LEVEL", "info"), "log level")
	flag.IntVar(&port, "port", utils.GetEnvInt("PORT", 8000), "port to expose faucet")

	flag.StringVar(&keyName, "key-name", utils.GetEnvString("KEY_NAME", faucet.DefaultKeyName), "the key name to be used by faucet")
	flag.StringVar(&keyMnemonic, "mnemonic", utils.GetEnvString("MNEMONIC", ""), "mnemonic for restoring key")
	flag.StringVar(&keyringPassword, "keyring-password", utils.GetEnvString("KEYRING_PASSWORD", ""), "the password for accessing keyring")
	flag.StringVar(&appCli, "cli-name", utils.GetEnvString("CLI_NAME", faucet.DefaultAppCli), "the name of the cli executable")
	flag.StringVar(&denom, "denom", utils.GetEnvString("DENOM", faucet.DefaultDenom), "the coin denomination")
	flag.Uint64Var(&creditAmount, "credit-amount", utils.GetEnvUint64("CREDIT_AMOUNT", faucet.DefaultCreditAmount), "the amount to credit in each request")
	flag.Uint64Var(&maxCredit, "max-credit", utils.GetEnvUint64("MAX_CREDIT", faucet.DefaultMaximumCredit), "the maximum credit per account")
}

func main() {
	flag.Parse()

	loggingLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(loggingLevel)

	f, err := faucet.NewFaucet(
		faucet.KeyName(keyName),
		faucet.Denom(denom),
		faucet.WithMnemonic(keyMnemonic),
		faucet.CliName(appCli),
		faucet.KeyringPassword(keyringPassword),
		faucet.CreditAmount(creditAmount),
		faucet.MaxCredit(maxCredit),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", f.ServeHTTP)
	log.Infof("listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
