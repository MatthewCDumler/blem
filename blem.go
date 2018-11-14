package main

import (
    "crypto/ecdsa"
    "context"
    "encoding/json"
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

type config struct {
    RawURL string
    UserKey string
    UserPW string
}

var userConfig config

// readConfig parses the json in the user's config for settings.
func readConfig() {
    jsonBlob, err := ioutil.ReadFile("/home/matt/go/src/github.com/MatthewCDumler/blem/blem.json")
    if err != nil {
            log.Fatal(err)
    }

    err = json.Unmarshal(jsonBlob, &userConfig)
    if err != nil {
            log.Fatal(err)
    }
}

// newPrivateKey returns a user's private key as a *ecdsa.PrivateKey.
func newPrivateKey(file, passwd string) (pk *ecdsa.PrivateKey) {
    // read in the json blob
    jsonBlob, err := ioutil.ReadFile(file)
    if err != nil {
            log.Fatal(err)
    }

    // decrypt the key
    key, err := keystore.DecryptKey(jsonBlob, passwd)
    if err != nil {
            log.Fatal(err)
    }

    // return the private key
    return key.PrivateKey
}

func main() {
    readConfig()
    // dial up client 
    client, err := ethclient.Dial(userConfig.RawURL)
    if err != nil {
            log.Fatal(err)
    }

    // get private key from keystore and password
    privateKey := newPrivateKey(userConfig.UserKey, userConfig.UserPW)
    // get public key from private key
    publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
    if !ok {
            log.Fatal("ECDSA ERROR")
    }

    // get user's address from public key
    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    // set the nonce
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
            log.Fatal(err)
    }
    // set the value, gas limit, and gas price
    value := big.NewInt(1)
    gasLimit := uint64(30000)
    gasPrice := big.NewInt(1)
    // get the receipient's address from hex
    toAddress := common.HexToAddress("0x14459c9a824599d0b98807729c9505a80b9e4f0f")

    // create the transaction
    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
    // get the chain ID
    chainID, err := client.NetworkID(context.Background())
    if err != nil {
            log.Fatal(err)
    }
    // sign the transaction
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
            log.Fatal(err)
    }

    // send the signed transaction
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
            log.Fatal(err)
    }
    // confirm transaction to output
    fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}
