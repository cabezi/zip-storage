package controller

import (
	"github.com/zipper-project/zipper/common/log"
	"github.com/zipper-project/zipper/common/rpc"
	"github.com/zipper-project/zipper/ledger/balance"
	pb "github.com/zipper-project/zipper/proto"
)

type chainData struct {
	cli *rpc.Client
}

func newChainData(rpcURL string) *chainData {
	client, err := rpc.DialHTTP(rpcURL)
	if err != nil {
		log.Error("create rpc client err: ", err)
	}
	return &chainData{cli: client}
}

func (c *chainData) getTx(hash string) (*pb.Transaction, error) {
	var tx *pb.Transaction
	if err := c.cli.Call("RPCLedger.GetTxByHash", hash, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *chainData) getBlock(height uint32) (*pb.Block, error) {
	var block *pb.Block
	if err := c.cli.Call("RPCLedger.GetBlockByHeight", height, block); err != nil {
		return nil, err
	}
	return block, nil
}

func (c *chainData) getBalance(addr string) (*balance.Balance, error) {
	var balance *balance.Balance
	if err := c.cli.Call("RPCLedger.GetBalance", addr, balance); err != nil {
		return nil, err
	}
	return balance, nil
}
