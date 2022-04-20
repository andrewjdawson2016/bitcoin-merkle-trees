package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ParsedBlock struct {
	Hash              string
	MerkleRoot        string
	TransactionHashes []string
}

const (
	HashKey        = "hash"
	MerkleRootKey  = "mrkl_root"
	TransactionKey = "tx"
)

func main() {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, Modules!")
		},
	}

	fmt.Println("Calling cmd.Execute()!")
	cmd.Execute()

	//blockId := "0000000000000000000013a86194dc6107ddc74dbabd3f2aa794d3722774cc03"
	//_, err := fetchBlock(blockId)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("hello world")
	//parsedBlock, err := parseBlock(block)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v", parsedBlock)
}

//func parseBlock(rawBlock []byte) (*ParsedBlock, error) {
//	m := make(map[string]interface{})
//	json.Unmarshal(rawBlock, &m)
//	blockHash, ok := m[HashKey]
//	if !ok {
//		return nil, errors.New("failed to find block hash")
//	}
//	blockHashStr, ok := blockHash.(string)
//	if !ok {
//		return nil, errors.New("failed to convert block has to string")
//	}
//	merkleRoot, ok := m[MerkleRootKey]
//	if !ok {
//		return nil, errors.New("failed to find merkle root")
//	}
//	merkleRootStr, ok := merkleRoot.(string)
//	if !ok {
//		return nil, errors.New("failed to convert merkle root to string")
//	}
//	transactions, ok := m[TransactionKey]
//	if !ok {
//		return nil, errors.New("failed to find transactions")
//	}
//
//}

func fetchBlock(blockId string) ([]byte, error) {
	url := fmt.Sprintf("https://blockchain.info/rawblock/%v", blockId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
