package table

import (
	"database/sql"
	"fmt"
	"time"
)

const TABLE_BALANCE = "BALANCE"

func init() {
	ac := NewBalance()
	Register(ac.TableName(), ac)
}

type Balance struct {
	ID      int64     `json:"id"`
	Addr    string    `json:"addr"`
	AssetID uint32    `json:"asset_id"`
	Amount  int64     `json:"amount"`
	Updated time.Time `json:"updated"`
}

//NewBalance return a tx object
func NewBalance() *Balance {
	return &Balance{}
}

//TableName return table name
func (balance *Balance) TableName() string {
	return TABLE_BALANCE
}

func (balance *Balance) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id BIGINT(20) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	addr VARCHAR(40) NOT NULL COMMENT '地址',
	asset_id BIGINT(20) NOT NULL COMMENT '资产ID',
	amount BIGINT(20) NOT NULL COMMENT '总额',
	updated DATETIME NOT NULL COMMENT '更新时间',	
	UNIQUE KEY uniq_addr_asset_id (addr,asset_id)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='资产表';`
	sql = fmt.Sprintf(sql, balance.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

func (balance *Balance) Insert(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("insert into %s(addr, asset_id, amount, updated) values(?, ?, ?, ?, ?, ?) on duplicate key update amount=?", balance.TableName()),
		balance.Addr, balance.AssetID, balance.Amount, balance.Updated, balance.Amount)
	if err != nil {
		return err
	}
	return nil
}
