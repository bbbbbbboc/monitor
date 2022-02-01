/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"flow/monitor/consts"
	"flow/monitor/utils"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/spf13/cobra"
)

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions <address>",
	Short: "Address' transactions watcher",
	Long:  `From now on, Start listing all transactions proposed / authorized / piad by the given address`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := flow.HexToAddress(args[0])
		if !address.IsValid(flow.Mainnet) {
			panic("invalid address")
		}
		// fmt.Println("start monitioring tx proposed / authorized / paid by address " + addressStr)
		monitorAddress(address)
	},
}

func init() {
	rootCmd.AddCommand(transactionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transactionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transactionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func monitorAddress(address flow.Address) {
	ctx := context.Background()
	flowClient := utils.NewFlowClient()

	// get the latest sealed block
	isSealed := true
	latestBlock, err := flowClient.GetLatestBlock(ctx, isSealed)
	utils.Handle(err)

	initialBlockHeight := latestBlock.BlockHeader.Height
	currentBlockHeight := initialBlockHeight
	for {
		currentBlock, err := flowClient.GetBlockByHeight(ctx, currentBlockHeight)
		// utils.PrintBlock(currentBlock)
		currentBlock = utils.HandleBlockDelayed(ctx, flowClient, currentBlockHeight, currentBlock, err)

		// check each collection
		for _, collectionGuarantee := range currentBlock.CollectionGuarantees {
			collectionId := collectionGuarantee.CollectionID
			collection, err := flowClient.GetCollection(ctx, collectionId)
			utils.Handle(err)

			// check each tx
			for _, transactionId := range collection.TransactionIDs {
				tx, err := flowClient.GetTransaction(ctx, transactionId)
				utils.Handle(err)

				// check address is payer / authorizer / proposer of current tx
				if utils.IsRelatedTransaction(tx, address) {
					utils.PrintTransaction(tx)
				}
			}
		}

		// flow block interval is rougly 3s according to block explorer, wait a while for next block
		time.Sleep(consts.BlockExistWaitInterval * time.Second)
		// go to next block
		currentBlockHeight++
	}
}
