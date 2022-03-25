package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

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

	switch sdkVersion {
	case "stargate-44":
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.StargateFortyFourVersion),
		)
	case "stargate-40":
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.StargateFortyVersion),
		)
	case "launchpad":
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.MaxLaunchpadVersion),
			chaincmd.WithLaunchpadCLI(appCli),
		)
	default:
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.Latest),
		)
	}

	cr, err := chaincmdrunner.New(context.Background(), chaincmd.New(appCli, ccoptions...))
	if err != nil {
		log.Fatal(err)
	}

	coins := strings.Split(defaultDenoms, denomSeparator)

	faucetOptions := make([]cosmosfaucet.Option, len(coins))
	for i, coin := range coins {
		faucetOptions[i] = cosmosfaucet.Coin(creditAmount, maxCredit, coin)
	}

	faucetOptions = append(faucetOptions, cosmosfaucet.Account(keyName, keyMnemonic, coinType))

	faucet, err := cosmosfaucet.New(context.Background(), cr, faucetOptions...)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", faucet.ServeHTTP)
	log.Infof("listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
