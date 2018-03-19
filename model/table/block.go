package table

import (
	"database/sql"
	"fmt"
	"time"
)

const TABLE_BLOCK = "block"

func init() {
	b := NewBlock()
	Register(b.TableName(), b)
}

type Block struct {
	ID            uint64    `json:"id"`
	Hash          string    `json:"hash"`
	PreviousHash  string    `json:"pre_hash"`
	StateHash     string    `json:"state_hash"`
	TimeStamp     time.Time `json:"created"`
	Nonce         uint32    `json:"nonce"`
	TxsMerkleHash string    `josn:"txs_merkle_hash"`
	Height        uint32    `json:"height"`
}

//NewBlock return a tx object
func NewBlock() *Block {
	return &Block{}
}

//TableName return table name
func (block *Block) TableName() string {
	return TABLE_BLOCK
}

func (block *Block) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	hash VARCHAR(255) NOT NULL COMMENT '区块哈希',
	pre_hash VARCHAR(255) NOT NULL COMMENT '前一个区块哈希',
	state_hash VARCHAR(255) NOT NULL COMMENT 'state哈希',
	created DATETIME NOT NULL COMMENT '创建时间',
	nonce BIGINT(20) UNSIGNED NOT NULL COMMENT '区块nonce',
	txs_merkle_hash VARCHAR(255) NOT NULL COMMENT '交易merkle哈希',	
	height BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '区块高度',
	UNIQUE KEY uniq_hash (hash)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='区块表';`
	sql = fmt.Sprintf(sql, block.TableName())
	_, err := db.Exec(sql)
	return sql, err
}
func (block *Block) QueryHeight(db *sql.DB) (uint32, error) {
	sql := fmt.Sprintf("select max(height) from %s", block.TableName())
	rows, err := db.Query(sql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var height uint32
	for rows.Next() {
		if err := rows.Scan(&height); err != nil {
			return 0, err
		}
	}
	return height, nil
}

func (block *Block) Insert(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("insert into%s(hash, pre_hash, state_hash, nonce, created, txs_merkle_hash, height) values(?, ?, ?, ?, ?, ?, ?)", block.TableName()),
		block.Hash, block.PreviousHash, block.StateHash, block.Nonce, block.TimeStamp, block.TxsMerkleHash, block.Height)
	if err != nil {
		return err
	}
	return nil
}
