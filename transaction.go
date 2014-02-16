package blockchain

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
