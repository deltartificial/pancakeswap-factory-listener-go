package main

import (
	"context"
	"fmt"
	"log"

	token "github.com/deltartificial/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/fatih/color"
)
  
var (
	  clientAddr   string = ""
	  pancakeFactory string = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
)

var routerAddress[2]string = [2]string{
	"0x10ED43C718714eb63d5aA57B78B54704E256024E", // PANCAKE SWAP
	"0xcF0feBd3f17CEf5b47b0cD257aCf6025c5BFf3b7"} // APE SWAP
	// ...

var pairTokensAddress [4]string = [4]string{
	"0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", // WBNB
	"0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d", // USDC
	"0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d", // BUSD
	"0x55d398326f99059fF775485246999027B3197955"} // USDT 

 

  func main() {

	  // Create a new client connections.
	  client, err := ethclient.Dial(clientAddr)
	  if err != nil {
		  log.Fatalf("error creating client: %s", err.Error())
	  }
	  defer client.Close()

  	  color.Set(color.FgYellow, color.Bold)
	  log.Print("BSC client connected!")
	  color.Unset() 
  
	  // Connect to ethereum contract.
	  tc, err := token.NewToken(common.HexToAddress(pancakeFactory), client)
	  if err != nil {
		  log.Fatalf("error connections to contract: %s", err.Error())
	  }
  
	  // Create pairCreated channel.
	  pairCreated := make(chan *token.TokenPairCreated)
  
	  // Subscribe to ethereum contract pairCreated event.
	  sub, err := tc.WatchPairCreated(nil, pairCreated, nil, nil)
	  if err != nil {
		  log.Fatalf("error subcribe to event: %s", err.Error())
	  }
  
	  for {
		  select {
		  case err := <-sub.Err():
			  log.Fatalf(": %s", err.Error())
		  case t := <-pairCreated:
			  log.Println("===================================================")
			  log.Printf("Token0:         %s", t.Token0)
			  log.Printf("Token1:           %s", t.Token1)
			  log.Printf("Pair:        %s", t.Pair)
			  log.Printf("Length:         %d", t.Arg3)
			  log.Printf("Block number: %d", t.Raw.BlockNumber)
			  log.Printf("Block hash:   %s", t.Raw.BlockHash)
			  log.Printf("Raw: 		%s", t.Raw.Data)
			  log.Printf("Address: 		%s", t.Raw.Address)
			  log.Printf("Topics: 		%s", t.Raw.Topics)
			  log.Printf("Transaction hash: %s", t.Raw.TxHash)
			  log.Printf("Transaction index: %d", t.Raw.TxIndex)
			  log.Printf("Log index: %d", t.Raw.Index)
			
			// ** RECOVER THE TOKEN TO SCAN
			token0String := t.Token0.String()
			token1String := t.Token1.String()

			getTokenToScan(token0String, token1String)
			
			txHash := t.Raw.TxHash
			fmt.Println("Transaction hash: ", txHash)
			log.Printf("Transaction hash: %s", txHash)

			tx, isPending, err := client.TransactionByHash(context.Background(), t.Raw.TxHash)
			if err != nil {
			 	  log.Fatal("TxHash ", err)
			}

			  fmt.Printf("TX Hash: %s\n", tx.Hash().Hex())
			  fmt.Printf("Pending?: %v\n", isPending)

			// Protected
			fmt.Printf("Protected?: %v\n", tx.Protected())
			// Type
			fmt.Printf("Type: %v\n", tx.Type())
			//   // ChainId
			fmt.Printf("ChainId: %v\n", tx.ChainId())
			//   // Gas Price
			fmt.Printf("Gas Price: %v\n", tx.GasPrice())
			//   // GasTipCap
			fmt.Printf("GasTipCap: %v\n", tx.GasTipCap())
			// To
			fmt.Printf("To: %v\n", tx.To())

			if tx.To() != nil {
				liquidity := tx.To().String()
				color.Set(color.FgGreen, color.Bold)
				fmt.Printf("ADD LIQUIDITY DETECTED: %s\n", liquidity)
				color.Unset()
				fmt.Println("===================================================")
			}

		  }
	  }
  }

  // ** RECOVER THE TOKEN TO SCAN
  func getTokenToScan(token0 string, token1 string) {
	for i := 0; i < len(pairTokensAddress); i++ {
		basePairTokenAddress := pairTokensAddress[i]
		if token0 == basePairTokenAddress {
			tokenToScan := token1
			fmt.Println("Token to scan: ", tokenToScan)
		} else if token1 == basePairTokenAddress {
			tokenToScan := token0
			fmt.Println("Token to scan: ", tokenToScan)
			
		}
	}
  }

  // ** CHECK IF THE ADDRESS "TO" IS A ROUTER
  func isRouter(to string) (bool) {
	for i := 0; i < len(routerAddress); i++ {
		if to == routerAddress[i] {
			return true
		} else {
			return false
		}
	}
	return false
  }
  