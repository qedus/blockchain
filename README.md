blockchain
==========

Blockchain.info API Go interface.

Check the test files for example use. However it is along the lines of:

```go

package main

import (
	"fmt"
	"net/http"

	"github.com/qedus/blockchain"
)

func main() {
	bc := blockchain.New(http.DefaultClient)
	address := &blockchain.Address{Address: "an address hash"}
	if err := bc.Request(address); err != nil {
		panic(err)
	}

	// Loop through all the transactions associated with a certain
	// address "an address hash" and print their transaction hashes.
	for {
		tx, err := address.NextTransaction()
		if err == blockchain.TransactionsDone {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Println(tx.Hash)
	}
}
```
