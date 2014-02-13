package blockchain

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	rootURL = "https://blockchain.info"
)

type BlockChain struct {
	client *http.Client
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

func (bc *BlockChain) httpGet(url string) (*http.Response, error) {
	resp, err := bc.client.Get(url)
	if err != nil {
		return nil, err
	}

	if err := checkHTTPResponse(resp); err != nil {
		return nil, err
	}
	return resp, nil
}
