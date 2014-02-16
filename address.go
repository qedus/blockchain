package blockchain

import (
	"errors"
	"fmt"
)

const (
	// Number of transactions to limit.
	// 50 is the maximum the API will allow.
	transactionLimit = 50
)

type Address struct {
	Hash160          string
	Address          string
	TransactionCount int64         `json:"n_tx"`
	TotalReceived    int64         `json:"total_received"`
	TotalSent        int64         `json:"total_sent"`
	FinalBalance     int64         `json:"final_balance"`
	Transactions     []Transaction `json:"txs"`

	// These are used for the NextTransaction iterator.
	bc         *BlockChain
	txOffset   int
	txPosition int
	txLimit    int
}

var TransactionsDone = errors.New("transactions done")

func (a *Address) NextTransaction() (*Transaction, error) {
	if a.txPosition < len(a.Transactions) {
		a.txPosition = a.txPosition + 1
		return &a.Transactions[a.txPosition-1], nil
	}

	if len(a.Transactions) < a.txLimit {
		return nil, TransactionsDone
	}
	if err := a.load(a.bc); err != nil {
		return nil, err
	}
	return a.NextTransaction()
}

func (a *Address) addressURL() string {
	// sort=1 orders transactions in ascending order
	return fmt.Sprintf(
		"%s/address/%s?format=json&sort=1&offset=%d&limit=%d",
		rootURL, a.Address, a.txOffset, a.txLimit)
}

func (a *Address) load(bc *BlockChain) error {
	a.bc = bc
	if a.txLimit == 0 {
		a.txLimit = transactionLimit
	}
	url := a.addressURL()
	if err := bc.httpGetJSON(url, a); err != nil {
		return err
	}
	a.txOffset = a.txOffset + transactionLimit
	a.txPosition = 0
	return nil
}
