package models

type AccountType string

const (
	CashAccount AccountType = "Cash"
	BankAccount AccountType = "Bank"
)

type Account struct {
	ID               string
	Type             AccountType
	Name             string
	Balance          float64
	LastTxnAmount    float64
	LastTxnTimestamp int64
}

type Event struct {
	ID        int64
	Message   string
	Timestamp int64
}