package client

type Command string

const (
	AddSavingsAccountCmd          Command = "CREATE_SAVINGS_ACCOUNT"
	SettleSavingsAccountCmd       Command = "SETTLE_SAVINGS_ACCOUNT"
	SearchTxnsByBankAccountCmd    Command = "SEARCH_BY_BANK_ACCOUNT"
	SearchTxnsBySavingsAccountCmd Command = "SEARCH_BY_SAVINGS_ACCOUNT"
)

func (cmd Command) ToString() string {
	return string(cmd)
}
