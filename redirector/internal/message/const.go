package message

type Command string
type Status string

const (
	AddSavingsAccountCmd          Command = "CREATE_ONLINE_SAVINGS_ACCOUNT"
	SettleSavingsAccountCmd       Command = "SETTLE_ONLINE_SAVINGS_ACCOUNT"
	SearchTxnsBySavingsAccountCmd Command = "SEARCH_SAVINGS_ACCOUNT"
)

func (cmd Command) ToString() string {
	return string(cmd)
}

const (
	SUCCESS Status = "success"
	ERROR   Status = "error"
)

func (s Status) ToString() string {
	return string(s)
}
