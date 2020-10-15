package faucet

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Faucet struct {
	appCli          string
	chainID         string
	keyringPassword string
	keyName         string
	faucetAddress   string
	keyMnemonic     string
	denom           string
	creditAmount    uint64
	maxCredit       uint64
}

func NewFaucet(opts ...Option) (*Faucet, error) {
	options := &defaultOptions
	for _, opt := range opts {
		opt(options)
	}

	e := Faucet{
		appCli:          options.AppCli,
		keyringPassword: options.KeyringPassword,
		keyName:         options.KeyName,
		keyMnemonic:     options.KeyMnemonic,
		denom:           options.Denom,
		creditAmount:    options.CreditAmount,
		maxCredit:       options.MaxCredit,
	}

	chainID, err := e.getChainID()
	if err != nil {
		return nil, err
	}
	e.chainID = chainID

	return &e, e.loadKey()
}

func (f *Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req httpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sent, err := f.GetTotalSent(req.Address)
		if err != nil {
			sendResponse(w, &httpResponse{
				Status: "failed",
				Error:  "could not get total tokens funded for this account",
			})
			return
		}

		if sent >= f.maxCredit {
			log.WithFields(map[string]interface{}{
				"address": req.Address,
				"amount":  fmt.Sprintf("%d%s", f.creditAmount, f.denom),
				"total":   sent + f.creditAmount,
			}).Warnf("tokens not sent: reached maximum credit")
			sendResponse(w, &httpResponse{
				Status: "failed",
				Error:  fmt.Sprintf("account has reached maximum credit allowed per account (%d)", f.maxCredit),
			})
			return
		}

		if err := f.Send(req.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.WithFields(map[string]interface{}{
			"address": req.Address,
			"amount":  fmt.Sprintf("%d%s", f.creditAmount, f.denom),
			"total":   sent + f.creditAmount,
		}).Infof("tokens sent")

		sendResponse(w, &httpResponse{
			Status: "ok",
		})
		return
	}
}

func sendResponse(w http.ResponseWriter, response *httpResponse) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
