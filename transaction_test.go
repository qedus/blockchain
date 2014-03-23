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

func TestTransactionFee(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	b := &blockchain.Block{Index: 312373}
	if err := bc.Request(b); err != nil {
		t.Fatal(err)
	}

	feeSum := int64(0)
	for _, tx := range b.Transactions {
		feeSum = feeSum + tx.Fee()
	}

	if feeSum != b.Fee {
		t.Fatalf("fees do not tally feeSum (%d) vs b.Fee (%d)",
			feeSum, b.Fee)
	}
}
