package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	chaincmdrunner "github.com/tendermint/starport/starport/pkg/chaincmd/runner"
	"github.com/tendermint/starport/starport/pkg/cosmosfaucet"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/xhttp"
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

	if sdkVersion == string(cosmosver.Stargate) {
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.StargateZeroFourtyAndAbove),
		)
	} else {
		ccoptions = append(ccoptions,
			chaincmd.WithVersion(cosmosver.LaunchpadAny),
			chaincmd.WithLaunchpadCLI(appCli),
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

	faucetOptions = append(faucetOptions, cosmosfaucet.Account(keyName, keyMnemonic))

	faucet, err := cosmosfaucet.New(context.Background(), cr, faucetOptions...)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", permitListMiddleware(faucet.ServeHTTP))
	log.Infof("listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func permitListMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// target index faucet (POST) handler in cosmosfaucet for permit list
		if r.URL.Path == "/" && r.Method == http.MethodPost {

			var req cosmosfaucet.TransferRequest

			// decode request into req.
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				transferResponseError(w, http.StatusBadRequest, err)
				return
			}

			if accountIsPermitted(req.AccountAddress) {
				// happy
				h.ServeHTTP(w, r) // call original handler
			} else {
				// not happy
				err := fmt.Errorf("%s is not permitted to receive a transfer from the faucet", req.AccountAddress)
				transferResponseError(w, http.StatusBadRequest, err)
				return
			}
		}
		h.ServeHTTP(w, r) // call original handler
	})
}

func transferResponseError(w http.ResponseWriter, code int, err error) {
	xhttp.ResponseJSON(w, code, cosmosfaucet.TransferResponse{
		Error: err.Error(),
	})
}
