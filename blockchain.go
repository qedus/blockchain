package blockchain

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	rootURL = "https://blockchain.info"
)

type BlockChain struct {
	client *http.Client
}

type Item interface {
	load(bc *BlockChain) error
}

func New(c *http.Client) *BlockChain {
	return &BlockChain{c}
}

func checkHTTPResponse(r *http.Response) error {
	if r.StatusCode == 200 {
		return nil
	}

	bodyErr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s: %.30q...", r.Status, bodyErr)
}

func (bc *BlockChain) Request(item Item) error {
	return item.load(bc)
}

func (bc *BlockChain) httpGetJSON(url string, v interface{}) error {
	resp, err := bc.client.Get(url)
	if err != nil {
		return err
	}

	if err := checkHTTPResponse(resp); err != nil {
		return err
	}

	defer resp.Body.Close()
	if err := decodeJSON(resp.Body, v); err != nil {
		return err
	}
	return nil
}

func decodeJSON(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)

	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}
