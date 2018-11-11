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

func main() {
	client, err := ethclient.Dial("enode://d6d5a3337f833adfed8de4c9a7ee0d6ddd8958a37af03a7df1348e2a7d32b0e3dd512c45c899245054dd2f110d4b1390f54745d5d2d95789414381d83359fb20@76.183.85.175:30303")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected")

	jsonBlob, err := ioutil.ReadFile("eth/.ethereum_private/keystore/UTC--2018-11-10T18-19-33.609052487Z--955f074178cdc5a7c72ba8fa0db9695fa9624446")
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey(jsonBlob, "apple123")
	if err != nil {
		log.Fatal(err)
	}
	pk := key.PrivateKey

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("ECDSA ERROR")
	}

	fAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
	balance, err := client.BalanceAt(context.Background(), fAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	nonce, err := client.PendingNonceAt(context.Background(), fAddr)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1)
	gasLimit := uint64(30000)
	gasPrice := big.NewInt(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gasLimit)
	fmt.Println(gasPrice)

	tAddr :=  common.HexToAddress("0xd7bfaddb6ea22e59aaefa337e0d182fdd44f6bf2")

	tx := types.NewTransaction(nonce, tAddr, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
