package table

import (
	"database/sql"
	"fmt"
	"time"
)

const TABLE_TRANSACTION = "transaction"

func init() {
	tx := NewTransaction()
	Register(tx.TableName(), tx)
}

type Transaction struct {
	ID                 uint64    `json:"id"`
	FromChain          string    `json:"from_chain"`
	ToChain            string    `json:"to_chain"`
	Type               uint32    `json:"type"`
	Nonce              uint32    `json:"nonce"`
	Sender             string    `json:"sender"`
	Receiver           string    `json:"receiver"`
	AssetID            uint32    `json:"asset_id"`
	Amount             int64     `json:"amount"`
	Fee                int64     `json:"fee"`
	Signature          string    `json:"signature"`
	Created            time.Time `json:"created"`
	Payload            string    `json:"payload"`
	Hash               string    `json:"hash"`
	Height             uint32    `json:"height"`
	ContractSpecAddr   string    `json:"contract_addr"`
	ContractSpecCode   string    `json:"contract_code"`
	ContractSpecParams string    `json:"contract_params"`
}

//NewTransaction return a tx object
func NewTransaction() *Transaction {
	return &Transaction{}
}

//TableName return table name
func (transaction *Transaction) TableName() string {
	return TABLE_TRANSACTION
}

func (transaction *Transaction) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	from_chain VARCHAR(255) NOT NULL COMMENT '发送链',
	to_chain VARCHAR(255) NOT NULL COMMENT '接受链',
	type TINYINT NOT NULL COMMENT '交易类型',
	nonce BIGINT(20) NOT NULL COMMENT '交易nonce',
	sender VARCHAR(255) NOT NULL COMMENT '发送地址',
	receiver VARCHAR(255) NOT NULL COMMENT '接收地址',
	asset_id BIGINT(20) NOT NULL COMMENT '资产ID',
	amount BIGINT(20) NOT NULL COMMENT '总额',
	fee BIGINT(20) NOT NULL COMMENT '手续费',
	signature VARCHAR(255) NOT NULL COMMENT '签名',
	created DATETIME NOT NULL COMMENT '创建时间',
	payload TEXT NOT NULL COMMENT '负载字段',
	hash VARCHAR(255) NOT NULL COMMENT '交易hash',
	height BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '所在区块高度',
	contract_addr VARCHAR(255) NOT NULL COMMENT '合约地址',
	contract_code TEXT NOT NULL COMMENT '合约代码',
	contract_params VARCHAR(255) NOT NULL COMMENT '合约参数',
	UNIQUE KEY uniq_hash (hash)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='交易表';`
	sql = fmt.Sprintf(sql, transaction.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

func (transaction *Transaction) Insert(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("insert into%s(from_chain, to_chain, type, nonce, sender, receiver, asset_id,amount,fee,signature,created,payload,hash,height,contract_addr,contract_code,contract_params) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", transaction.TableName()),
		transaction.FromChain,
		transaction.ToChain,
		transaction.Type,
		transaction.Nonce,
		transaction.Sender,
		transaction.Receiver,
		transaction.AssetID,
		transaction.Amount,
		transaction.Fee,
		transaction.Signature,
		transaction.Created,
		transaction.Payload,
		transaction.Hash,
		transaction.Height,
		transaction.ContractSpecAddr,
		transaction.ContractSpecCode,
		transaction.ContractSpecParams)
	if err != nil {
		return err
	}
	return nil
}
