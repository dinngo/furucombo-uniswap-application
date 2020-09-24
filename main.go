package main

import (
	_ "github.com/joho/godotenv/autoload"

	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	uniswapAPI       = "https://gentle-frost-9e74.uniswap.workers.dev/1/"
	txsFileName      = "txs.txt"
	accountsFileName = "accounts.txt"
)

func main() {
	endpoint := os.Getenv("ENDPOINT")

	senders := make([]common.Address, 0)
	err := getSendersFromTxs(endpoint, &senders)
	if err != nil {
		panic(err)
	}

	accounts := make([]string, 0)
	if err := filterOutSenders(&senders, &accounts); err != nil {
		panic(err)
	}

	if err = outputAccounts(&accounts); err != nil {
		panic(err)
	}
}

func getSendersFromTxs(endpoint string, senders *[]common.Address) error {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return err
	}

	f, err := os.Open(txsFileName)
	if err != nil {
		return err
	}

	checkExisting := make(map[common.Address]struct{}, 0)
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

		if _, ok := checkExisting[from]; !ok {
			checkExisting[from] = struct{}{}
			// Fill in senders
			*senders = append(*senders, from)
			fmt.Println("from: ", from.Hex(), hash.Hex())
		}
	}

	return nil
}

func filterOutSenders(senders *[]common.Address, accounts *[]string) error {
	for _, sender := range *senders {
		resp, err := http.Get(uniswapAPI + sender.Hex())
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			*accounts = append(*accounts, sender.Hex())
			fmt.Println("address: ", sender.Hex())
		}
	}

	return nil
}

func outputAccounts(accounts *[]string) error {
	f, err := os.Create(accountsFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range *accounts {
		_, err := f.WriteString(v + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
