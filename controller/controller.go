package controller

import (
	"database/sql"
	"strings"
	"time"

	"github.com/cabezi/zip-storage/config"
	"github.com/cabezi/zip-storage/model"
	"github.com/cabezi/zip-storage/model/table"
	"github.com/zipper-project/zipper/common/log"
	pb "github.com/zipper-project/zipper/proto"
)

const contractParamSeparator = "$"

type Controller struct {
	tmpHeight uint32
	cd        *chainData
	no        *notify
}

func NewControler(cfg *config.Config) *Controller {
	mysqlHeight, err := table.NewBlock().QueryHeight(model.DB)
	if err != nil {
		log.Error(err)
	}
	return &Controller{
		cd:        newChainData(cfg.RPCURL),
		no:        newNotify(cfg.WSURL, cfg.ORIGIN),
		tmpHeight: mysqlHeight,
	}
}

func (c *Controller) ReceivePushHeight() {
	for {
		select {
		case height := <-c.no.heightChannel():
			log.Debugln("recive websocket push height ", height)
			for i := c.tmpHeight; i < c.checkBlockHeight(height); i++ {
				block, err := c.cd.getBlock(i)
				if err != nil {
					log.Error("getblock err:", err)
				}

				tx, err := model.DB.Begin()
				if err != nil {
					log.Error(err)
				}

				if err := c.putBlockHeader(block.Header, tx); err != nil {
					log.Error(err)
					tx.Rollback()
				}

				if err := c.putTransactions(block.Transactions, block.Height(), tx); err != nil {
					log.Error(err)
					tx.Rollback()
				}
				tx.Commit()
			}
		}
	}
}

func (c *Controller) checkBlockHeight(height uint32) uint32 {
	if c.tmpHeight < height {
		return height - c.tmpHeight
	}
	panic("err height")
}

func (c *Controller) putBlockHeader(block *pb.BlockHeader, stx *sql.Tx) error {
	b := &table.Block{
		Hash:          block.Hash().String(),
		PreviousHash:  block.PreviousHash,
		StateHash:     block.StateHash,
		TimeStamp:     time.Unix(int64(block.TimeStamp), 0),
		Nonce:         block.Nonce,
		TxsMerkleHash: block.TxsMerkleHash,
		Height:        block.Height,
	}
	return b.Insert(stx)
}

func (c *Controller) putTransactions(txs pb.Transactions, height uint32, stx *sql.Tx) error {
	for _, tx := range txs {
		ttx := &table.Transaction{
			FromChain: string(tx.Header.FromChain),
			ToChain:   string(tx.Header.ToChain),
			Type:      uint32(tx.Header.Type),
			Nonce:     tx.Header.Nonce,
			Sender:    tx.Header.Sender,
			Receiver:  tx.Header.Recipient,
			AssetID:   tx.Header.AssetID,
			Amount:    tx.Header.Amount,
			Fee:       tx.Header.Fee,
			Signature: string(tx.Header.Signature),
			Created:   time.Unix(int64(tx.Header.CreateTime), 0),
			Payload:   string(tx.Payload),
			Hash:      tx.Hash().String(),
			Height:    height,
		}

		if tx.ContractSpec != nil {
			ttx.ContractSpecAddr = string(tx.ContractSpec.Addr)
			ttx.ContractSpecCode = string(tx.ContractSpec.Code)
			ttx.ContractSpecParams = strings.Join(tx.ContractSpec.Params, contractParamSeparator)
		}

		if err := ttx.Insert(stx); err != nil {
			return err
		}

		balance, err := c.cd.getBalance(ttx.Sender)
		if err != nil {
			return err
		}

		if err := c.putBalanace(ttx.AssetID, balance.Amounts[ttx.AssetID], ttx.Sender, stx); err != nil {
			return err
		}

		balance, err = c.cd.getBalance(ttx.Receiver)
		if err != nil {
			return err
		}

		if err := c.putBalanace(ttx.AssetID, balance.Amounts[ttx.AssetID], ttx.Receiver, stx); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) putBalanace(assetID uint32, amount int64, addr string, stx *sql.Tx) error {
	b := &table.Balance{
		Addr:    addr,
		AssetID: assetID,
		Amount:  amount,
		Updated: time.Now(),
	}
	return b.Insert(stx)
}
