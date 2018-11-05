/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"time"

	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/cobra"
	"github.com/Oneledger/protocol/node/chains/bitcoin"

	"github.com/btcsuite/btcd/chaincfg"
)

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait for something to happen",
	Run:   Wait,
}

func init() {
	RootCmd.AddCommand(waitCmd)

	var completed bool
	var strings []string

	waitCmd.Flags().BoolVar(&completed, "completed", false, "send recipient")
	waitCmd.Flags().StringArrayVar(&strings, "party", strings, "send recipient")
}

// TODO: Wait for real things to happen...
func Wait(cmd *cobra.Command, args []string) {
	log.Debug("Waiting")
	cli := bitcoin.GetBtcClient("127.0.0.1:18833", &chaincfg.RegressionNetParams)
	stop := bitcoin.ScheduleBlockGeneration(*cli, 10 )
	time.Sleep(60 * time.Second)
	bitcoin.StopBlockGeneration(stop)
}
