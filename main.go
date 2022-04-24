package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var (
	blockId         string
	numTransactions int
)

type ParsedBlock struct {
	Hash              string
	MerkleRoot        string
	TransactionHashes []string
}

const (
	RawBlockURITemplate   = "https://blockchain.info/rawblock/%v"
	BlockHashPath         = "hash"
	MerkleRootPath        = "mrkl_root"
	TransactionHashesPath = "tx.#.hash"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mrklctl",
		Short: "A CLI tool for downloading blocks and understanding their Merkle Trees",
	}
	getBlockCmd = &cobra.Command{
		Use:   "get_block",
		Short: "Fetches a block with given ID and outputs the parts the block used in the Merkle tree construction",
		Run: func(cmd *cobra.Command, args []string) {
			if err := getBlock(blockId, numTransactions); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	getBlockCmd.PersistentFlags().StringVarP(&blockId, "block_id", "", "", "Block of targeted block")
	getBlockCmd.PersistentFlags().IntVarP(&numTransactions, "num_transactions", "", 10, "Number of transactions to print, 0 for all transactions")
}

func main() {
	rootCmd.AddCommand(getBlockCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func getBlock(blockId string, numTransactions int) error {
	if blockId == "" {
		return errors.New("block_id is required")
	}
	rawBlock, err := fetchBlock(blockId)
	if err != nil {
		return err
	}
	pb, err := parseBlock(rawBlock)
	if err != nil {
		return err
	}
	printBlock(pb)
	return nil
}

func printBlock(pb *ParsedBlock) {
	fmt.Println("Block Hash: ", pb.Hash)
	fmt.Println("Merkle Root: ", pb.MerkleRoot)
	fmt.Println("=== Transactions ===")
	lastIndex := numTransactions - 1
	if numTransactions == 0 || lastIndex > len(pb.TransactionHashes)-1 {
		lastIndex = len(pb.TransactionHashes) - 1
	}
	for i := 0; i <= lastIndex; i++ {
		fmt.Println(pb.TransactionHashes[i])
	}
}

func parseBlock(rawBlock string) (*ParsedBlock, error) {
	result := &ParsedBlock{}
	hashKey := gjson.Get(rawBlock, BlockHashPath)
	if !hashKey.Exists() {
		return nil, fmt.Errorf("failed to extract %v from block", BlockHashPath)
	}
	result.Hash = hashKey.String()

	merkleRoot := gjson.Get(rawBlock, MerkleRootPath)
	if !merkleRoot.Exists() {
		return nil, fmt.Errorf("failed to extract %v from block", MerkleRootPath)
	}
	result.MerkleRoot = merkleRoot.String()
	
	transactionHashes := gjson.Get(rawBlock, TransactionHashesPath)
	if !transactionHashes.Exists() || !transactionHashes.IsArray() {
		return nil, fmt.Errorf("failed to extract %v from block", TransactionHashesPath)
	}
	transactionHashes.ForEach(func(_, value gjson.Result) bool {
		result.TransactionHashes = append(result.TransactionHashes, value.String())
		return true
	})
	return result, nil
}

func fetchBlock(blockId string) (string, error) {
	url := fmt.Sprintf(RawBlockURITemplate, blockId)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
