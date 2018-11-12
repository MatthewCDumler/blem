package main

import (
	"crypto/ecdsa"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	privateKey *ecdsa.PrivateKey
)

func privateKeyInit(file, passwd string) {
	jsonBlob, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	key, err := keystore.DecryptKey(jsonBlob, passwd)
	if err != nil {
		log.Fatal(err)
	}
	privateKey = key.PrivateKey
}

func main() {
	// dial up client 
	client, err := ethclient.Dial("ethereum/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyInit("ethereum/keystore/UTC--2018-11-11T00-54-48.985136186Z--83a184bcc727851c493e6a2a76bcd37b86245de6", "apple")
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("ECDSA ERROR")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(1)
	gasLimit := uint64(30000)
	gasPrice := big.NewInt(1)
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0x14459c9a824599d0b98807729c9505a80b9e4f0f")

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}
