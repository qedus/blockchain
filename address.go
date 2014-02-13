package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// Number of transactions to limit.
	// 50 is the maximum the API will allow.
	transactionLimit = 50
)

type Input struct {
	PrevOut struct {
		Address          string `json:"addr"`
		Number           int64  `json:"n"`
		TransactionIndex int64  `json:"tx_index"`
		Type             int64
		Value            int64
	} `json:"prev_out"`
}

type Output struct {
	Address          string `json:"addr"`
	AddressTag       string `json:"addr_tag"`
	AddressTagLink   string `json:"addr_tag_link"`
	Number           int64  `json:"n"`
	TransactionIndex int64  `json:"tx_index"`
	Type             int64
	Value            int64
}

type Transaction struct {
	Hash             string
	Inputs           []Input
	InputCount       int64    `json:"vin_sz"`
	Outputs          []Output `json:"out"`
	OutputCount      int64    `json:"vout_sz"`
	RelayedBy        string   `json:"relayed_by"`
	Result           int64
	Size             int64
	Time             int64
	BlockHeight      int64 `json:"block_height"`
	TransactionIndex int64 `json:"tx_index"`
	Version          int64 `json:"ver"`
}

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
	hash       string
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
	if err := a.bc.httpGetAddress(a); err != nil {
		return nil, err
	}
	return a.NextTransaction()
}

func addressURL(hash string, offset, limit int) string {
	// sort=1 orders transactions in ascending order
	return fmt.Sprintf(
		"%s/address/%s?format=json&sort=1&offset=%d&limit=%d",
		rootURL, hash, offset, limit)
}

func (bc *BlockChain) GetAddress(hash string) (*Address, error) {
	address := &Address{bc: bc, hash: hash,
		txOffset: 0, txLimit: transactionLimit}
	return address, bc.httpGetAddress(address)
}

func (bc *BlockChain) httpGetAddress(address *Address) error {
	url := addressURL(address.hash, address.txOffset, address.txLimit)
	resp, err := bc.httpGet(url)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	if err := dec.Decode(address); err != nil {
		return err
	}

	address.txOffset = address.txOffset + transactionLimit
	address.txPosition = 0
	return nil
}
