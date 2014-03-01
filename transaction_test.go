package blockchain_test

import (
	"github.com/qedus/blockchain"
	"net/http"
	"testing"
)

func TestUnconfirmedTransactions(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	ut := &blockchain.UnconfirmedTransactions{}
	if err := bc.Request(ut); err != nil {
		t.Fatal(err)
	}

	if len(ut.Transactions) == 0 {
		t.Fatal("no transactions")
	}

	count := 0
	for {

		tx, err := ut.NextTransaction()

		if err == blockchain.IterDone {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		if tx.Hash == "" {
			t.Fatal("no transaction hash")
		}
		count++
	}
	t.Logf("%d unconfirmed transactions", count)
}
