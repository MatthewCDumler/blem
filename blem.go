package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("enode://d33e58a0320d9d0a07207b62c80073b6bc41b2735cba80b64bc403994d0862552b52452b87be93ef4736e8be3f9fe2e18fff9eddda0946181a8a51f710200534@76.183.85.175:30303")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected")

	account := common.HexToAddress("0x955f074178cdc5a7c72ba8fa0db9695fa9624446")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)
}
