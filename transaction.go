package blockchain

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	// The maximum number of transactions that can be requested in one
	// API call.
	maxTransactionLimit = 50
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

func (t *Transaction) IsCoinbase() bool {
	return len(t.Inputs) == 0
}

func (t *Transaction) Fee() int64 {
	if t.IsCoinbase() {
		return 0
	}

	inputSum := int64(0)
	for _, input := range t.Inputs {
		inputSum = inputSum + input.PrevOut.Value
	}

	outputSum := int64(0)
	for _, output := range t.Outputs {
		outputSum = outputSum + output.Value
	}
	return inputSum - outputSum
}

type UnconfirmedTransactions struct {
	Transactions []Transaction `json:"txs"`

	// These are used for NextTransaction iterator.
	bc         *BlockChain
	txOffset   int
	txPosition int
	txLimit    int
}

func (uc *UnconfirmedTransactions) NextTransaction() (*Transaction, error) {
	if uc.txPosition < len(uc.Transactions) {
		uc.txPosition = uc.txPosition + 1
		return &uc.Transactions[uc.txPosition-1], nil
	}

	if len(uc.Transactions) < uc.txLimit {
		return nil, IterDone
	}
	if err := uc.load(uc.bc); err != nil {
		return nil, err
	}
	return uc.NextTransaction()
}

func (ut *UnconfirmedTransactions) unconfirmedTransactionsURL() string {
	v := url.Values{}
	v.Set("format", "json")
	v.Set("sort", "0")
	v.Set("offset", strconv.Itoa(ut.txOffset))
	v.Set("limit", strconv.Itoa(ut.txLimit))
	return fmt.Sprintf("%s/unconfirmed-transactions?%s",
		rootURL, v.Encode())
}

func (ut *UnconfirmedTransactions) load(bc *BlockChain) error {
	ut.bc = bc
	if ut.txLimit == 0 {
		ut.txLimit = maxTransactionLimit
	}
	url := ut.unconfirmedTransactionsURL()
	if err := bc.httpGetJSON(url, ut); err != nil {
		return err
	}
	ut.txOffset = ut.txOffset + ut.txLimit
	ut.txPosition = 0
	return nil
}
