package model

import "time"

type HistoryDrug struct {
	Drug      *Drug     `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}
