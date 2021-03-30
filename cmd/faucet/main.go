package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/cosmosver"

	"github.com/tendermint/faucet/internal/environ"
)

var (
	port            int
	keyringBackend  string
	sdkVersion      string
	keyName         string
	keyMnemonic     string
	keyringPassword string
	appCli          string
	denom           string
	creditAmount    uint64
	maxCredit       uint64
	nodeAddress     string
	legacySendCmd   bool
)

func init() {
	flag.IntVar(&port, "port", environ.GetInt("PORT", 8000), "port to expose faucet")
	flag.StringVar(&keyringBackend, "keyring-backend", environ.GetString("KEYRING_BACKEND", ""), "keyring backend")
	flag.StringVar(&sdkVersion, "sdk-version", environ.GetString("SDK_VERSION", "stargate"), "version of sdk (launchpad or stargate)")
	flag.StringVar(&keyName, "key-name", environ.GetString("KEY_NAME", cosmosfaucet.DefaultAccountName), "the key name to be used by faucet")
	flag.StringVar(&keyMnemonic, "mnemonic", environ.GetString("MNEMONIC", ""), "mnemonic for restoring key")
	flag.StringVar(&keyringPassword, "keyring-password", environ.GetString("KEYRING_PASSWORD", ""), "the password for accessing keyring")
	flag.StringVar(&appCli, "cli-name", environ.GetString("CLI_NAME", "gaiad"), "the name of the cli executable")
	flag.StringVar(&denom, "denom", environ.GetString("DENOM", cosmosfaucet.DefaultDenom), "the coin denomination")
	flag.Uint64Var(&creditAmount, "credit-amount", environ.GetUint64("CREDIT_AMOUNT", cosmosfaucet.DefaultAmount), "the amount to credit in each request")
	flag.Uint64Var(&maxCredit, "max-credit", environ.GetUint64("MAX_CREDIT", cosmosfaucet.DefaultMaxAmount), "the maximum credit per account")
	flag.StringVar(&nodeAddress, "node", environ.GetString("NODE", ""), "the address of the node that will handle the requests")
	flag.BoolVar(&legacySendCmd, "legacy-send", environ.GetBool("LEGACY_SEND", false), "use legacy send command")
}

func main() {
	flag.Parse()

	configKeyringBackend, err := chaincmd.KeyringBackendFromString(keyringBackend)
	if err != nil {
		log.Fatal(err)
	}

	ccoptions := []chaincmd.Option{
		chaincmd.WithKeyringPassword(keyringPassword),
		chaincmd.WithKeyringBackend(configKeyringBackend),
		chaincmd.WithAutoChainIDDetection(),
		chaincmd.WithNodeAddress(nodeAddress),
	}

	if legacySendCmd {
		ccoptions = append(ccoptions, chaincmd.WithLegacySendCommand())
	}

	if sdkVersion == "stargate" {
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.StargateZeroFourtyAndAbove),
		)
	} else {
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.LaunchpadAny),
			chaincmd.WithLaunchpadCLI(appCli),
		)
	}

	cc := chaincmd.New(appCli, ccoptions...)
	cr, err := chaincmdrunner.New(context.Background(), cc)
	if err != nil {
		log.Fatal(err)
	}

	faucet, err := cosmosfaucet.New(
		context.Background(),
		cr,
		cosmosfaucet.Account(keyName, keyMnemonic),
		cosmosfaucet.Coin(creditAmount, maxCredit, denom),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", faucet.ServeHTTP)
	log.Infof("listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
