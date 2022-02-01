package utils

import (
	"context"
	"flow/monitor/consts"
	"fmt"
	"strings"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func NewFlowClient() *client.Client {
	flowClient, err := client.New(consts.AccessAPI, grpc.WithInsecure())
	Handle(err)

	return flowClient
}

func Handle(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

// Keep trying to fetch current block if not exist
func HandleBlockDelayed(
	ctx context.Context, flowClient *client.Client, currentBlockHeight uint64, currentBlock *flow.Block, err error,
) *flow.Block {
	for {
		if err == nil {
			return currentBlock
		}
		if strings.Contains(err.Error(), "NotFound") {
			fmt.Printf(
				"block has not been produced yet, keep trying to fetch it every %d seconds, ctrl c to stop\n",
				consts.BlockNotExistWaitInterval,
			)
			time.Sleep(consts.BlockNotExistWaitInterval * time.Second)
			currentBlock, err = flowClient.GetBlockByHeight(ctx, currentBlockHeight)
		} else {
			Handle(err)
		}
	}
}

// Return true if address is authorizer / payer / proposer of the tx
func IsRelatedTransaction(tx *flow.Transaction, address flow.Address) bool {
	if tx.Payer == address || tx.ProposalKey.Address == address || sliceContains(tx.Authorizers, address) {
		return true
	}
	return false
}

func PrintBlock(block *flow.Block) {
	fmt.Printf("\nblock ID: %s\n", block.ID)
	fmt.Printf("block height: %d\n", block.Height)
	fmt.Printf("block timestamp: %s\n", block.Timestamp)
}

func PrintTransaction(tx *flow.Transaction) {
	fmt.Println(tx.ID().String())
	// fmt.Printf("tx ID: %s", tx.ID().String())
	// fmt.Printf("tx Payer: %s\n", tx.Payer.String())
	// fmt.Printf("tx Proposer: %s\n", tx.ProposalKey.Address.String())
	// fmt.Printf("tx Authorizers: %s\n\n", tx.Authorizers)
}

func sliceContains(authorizers []flow.Address, address flow.Address) bool {
	for _, authorizer := range authorizers {
		if authorizer == address {
			return true
		}
	}
	return false
}
