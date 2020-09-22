package main

import (
	_ "github.com/joho/godotenv/autoload"

	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	uniswapAPI   = "https://gentle-frost-9e74.uniswap.workers.dev/1/"
	txsFileName  = "txs.txt"
	jsonFileName = "affected_addresses.json"
)

// Report represents affected addresses
type Report struct {
	Addresses []string `json:"addresses"`
}

func main() {
	endpoint := os.Getenv("ENDPOINT")

	senders := make(map[common.Address]struct{})
	err := getSendersFromTxs(endpoint, senders)
	if err != nil {
		panic(err)
	}

	report := new(Report)
	if err := filterOutSenders(senders, report); err != nil {
		panic(err)
	}

	if err = outputJSON(report); err != nil {
		panic(err)
	}
}

func getSendersFromTxs(endpoint string, senders map[common.Address]struct{}) error {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return err
	}

	// txs.txt is derived from https://explore.duneanalytics.com/public/dashboards/Cw10tmbxQ09KemdZdbzkW9rFMWTsBdE3CWW27jIF
	f, err := os.Open(txsFileName)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hash := common.HexToHash(scanner.Text())
		tx, isPending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil {
			return err
		}
		if isPending {
			return errors.New("unexpected isPending")
		}

		msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
		if err != nil {
			return err
		}
		from := msg.From()

		if _, ok := senders[from]; !ok {
			senders[from] = struct{}{}
			fmt.Println("from: ", from.Hex(), hash.Hex())
		}
	}

	return nil
}

func filterOutSenders(senders map[common.Address]struct{}, report *Report) error {
	for sender := range senders {
		resp, err := http.Get(uniswapAPI + sender.Hex())
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			report.Addresses = append(report.Addresses, sender.Hex())
			fmt.Println("address: ", sender.Hex())
		}
	}

	return nil
}

func outputJSON(report *Report) error {
	jsonBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil
	}

	if err := ioutil.WriteFile(jsonFileName, jsonBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}
