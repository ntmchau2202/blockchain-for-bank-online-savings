package endpoint

import (
	"gr-blockchain-side/internal/message"
)

func CreateErrorMessage(msg string) (resp message.Response) {
	resp.Stat = message.ERROR
	resp.Details = make(map[string]interface{})
	resp.Details["message"] = msg
	return
}

func CreateSuccessMessage(msg string, txn string, blockNumber string) (resp message.Response) {
	resp.Stat = message.SUCCESS
	resp.Details = make(map[string]interface{})
	resp.Details["message"] = msg
	resp.Details["transaction_id"] = txn
	resp.Details["block_number"] = blockNumber
	return
}
