package blockchain_test

import (
	"github.com/qedus/blockchain"
	"net/http"
	"testing"
)

const (
	largeAddress               = "1dice8EMZmqKvrGE4Qc9bUFf9PX3xaYDp"
	largeAddressTxHashOne      = "a0a87b1577d606b349cfded85c842bdc53b99bcd49614229a71804b46b1c27cc"
	largeAddressTxHashFiftyOne = "195ad7c67887bb4b03cea7398f70dc36b9574c80a2d46382e557000a52f1a65e"

	smallAddress = "1416ArzSr5HGzaTbfQHjkLE5RVBBGw3W13"
)

func TestGetLargeAddress(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	address, err := bc.GetAddress(largeAddress)
	if err != nil {
		t.Fatal(err)
	}

	if len(address.Transactions) != 50 {
		t.Fatalf("tx count not 50")
	}

	tx, err := address.NextTransaction()
	if err != nil {
		t.Fatal(err)
	}
	if tx.Hash != largeAddressTxHashOne {
		t.Fatalf("tx hash incorrect %s vs %s",
			tx.Hash, largeAddressTxHashOne)
	}

	for i := 0; i < 49; i++ {
		tx, err = address.NextTransaction()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Check the address iterator goes to the server again.
	tx, err = address.NextTransaction()
	if err != nil {
		t.Fatal(err)
	}
	if tx.Hash != largeAddressTxHashFiftyOne {
		t.Fatalf("tx hash incorrect %s vs %s",
			tx.Hash, largeAddressTxHashFiftyOne)
	}
	if len(address.Transactions) != 50 {
		t.Fatalf("tx count not 50")
	}
}

func TestGetSmallAddress(t *testing.T) {
	bc := blockchain.New(http.DefaultClient)
	address, err := bc.GetAddress(smallAddress)
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for {
		_, err := address.NextTransaction()
		if err == blockchain.TransactionsDone {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		count++
	}

	if count != 6 {
		t.Fatalf("expected 6 iterations but got %d", count)
	}
}
