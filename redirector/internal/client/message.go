package client

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout string = "Mon, 02 Jan 2006 15:04:05 MST"
)

type Message struct {
	Command    string `json:"command"`
	TimeIssued string `json:"time_issued"`
}

func validateCommand(requestedCmd string) (err error) {
	if strings.Compare(requestedCmd, AddSavingsAccountCmd.ToString()) != 0 && strings.Compare(requestedCmd, SettleSavingsAccountCmd.ToString()) != 0 && strings.Compare(requestedCmd, SearchTxnsByBankAccountCmd.ToString()) != 0 && strings.Compare(requestedCmd, SearchTxnsBySavingsAccountCmd.ToString()) != 0 {
		return errors.New("undefined command")
	}
	return nil
}

func validateTime(timeString string) (err error) {
	if _, err = time.Parse(timeLayout, timeString); err != nil {
		return errors.New("invalid time format")
	}
	return nil
}

func (m Message) Validate() (err error) {

	if err = validateTime(m.TimeIssued); err != nil {
		return
	}

	return validateCommand(m.Command)
}

type AddSavingsAccountMessage struct {
	Message
	Details struct {
		BankAccountNumber    string `json:"bank_acc_number"`
		SavingsAccountNumber string `json:"savings_acc_number"`
		TimeCreated          string `json:"time_created"`
		InitialAmount        string `json:"initial_amount"`
		InterestRate         string `json:"interest_rate"` // in float
		TypeOfSavings        string `json:"type_of_savings"`
		SavingsPeriod        string `json:"saving_period"` // in months
		TransactionUnit      string `json:"transaction_unit"`
	} `json:"details"`
}

func (m AddSavingsAccountMessage) Validate() (err error) {

	if err = m.Message.Validate(); err != nil {
		return
	}

	if strings.Compare(m.Message.Command, AddSavingsAccountCmd.ToString()) != 0 {
		return errors.New("invalid command")
	}

	if len(m.Details.BankAccountNumber) == 0 {
		return errors.New("missing bank account number")
	}

	if len(m.Details.SavingsAccountNumber) == 0 {
		return errors.New("missing saving account number")
	}

	if err = validateTime(m.Details.TimeCreated); err != nil {
		return
	}

	if _, err := strconv.ParseFloat(m.Details.InitialAmount, 64); err != nil {
		return errors.New("invalid initial savings amount")
	}

	if _, err := strconv.ParseFloat(m.Details.InterestRate, 64); err != nil {
		return errors.New("invalid interest rate")
	}

	if _, err := strconv.ParseInt(m.Details.SavingsPeriod, 10, 64); err != nil {
		return errors.New("invalid saving period")
	}

	if len(m.Details.TransactionUnit) != 3 || strings.Compare(strings.ToUpper(m.Details.TransactionUnit), m.Details.TransactionUnit) != 0 {
		return errors.New("invalid transaction unit format")
	}
	return nil
}

type SettleSavingsAccountMessage struct {
	Message
	Details struct {
		BankAccountNumber    string `json:"bank_acc_number"`
		SavingsAccountNumber string `json:"savings_acc_number"`
		TimeSettled          string `json:"time_settled"`
		InterestAmount       string `json:"interst_amount"`
		TotalAmount          string `json:"total_amount"` // total balance = initial + interest
		TransactionUnit      string `json:"transaction_unit"`
	} `json:"details"`
}

func (m SettleSavingsAccountMessage) Validate() (err error) {
	if err = m.Message.Validate(); err != nil {
		return
	}

	if m.Details == struct {
		BankAccountNumber    string "json:\"bank_acc_number\""
		SavingsAccountNumber string "json:\"savings_acc_number\""
		TimeSettled          string "json:\"time_settled\""
		InterestAmount       string "json:\"interst_amount\""
		TotalAmount          string "json:\"total_amount\""
		TransactionUnit      string "json:\"transaction_unit\""
	}{} {
		return errors.New("missing details")
	}

	if strings.Compare(m.Message.Command, SettleSavingsAccountCmd.ToString()) != 0 {
		return errors.New("invalid command")
	}

	if len(m.Details.BankAccountNumber) == 0 {
		return errors.New("missing bank account number")
	}

	if len(m.Details.SavingsAccountNumber) == 0 {
		return errors.New("missing saving account number")
	}

	if err = validateTime(m.Details.TimeSettled); err != nil {
		return
	}

	if _, err := strconv.ParseFloat(m.Details.InterestAmount, 64); err != nil {
		return errors.New("invalid interest amount")
	}

	if _, err := strconv.ParseFloat(m.Details.TotalAmount, 64); err != nil {
		return errors.New("invalid total returned amount")
	}

	if len(m.Details.TransactionUnit) != 3 || strings.Compare(strings.ToUpper(m.Details.TransactionUnit), m.Details.TransactionUnit) != 0 {
		return errors.New("invalid transaction unit format")
	}
	return nil
}

type QueryMessage struct {
	Message
	Details struct {
		QueryID string `json:"query_id"`
	} `json:"details"`
}

func (m QueryMessage) Validate() (err error) {
	if err = m.Message.Validate(); err != nil {
		return
	}

	if strings.Compare(m.Message.Command, SearchTxnsByBankAccountCmd.ToString()) != 0 && strings.Compare(m.Message.Command, SearchTxnsBySavingsAccountCmd.ToString()) != 0 {
		return errors.New("invalid command")
	}

	if m.Details == struct {
		QueryID string "json:\"query_id\""
	}{} {
		return errors.New("missing details")
	}

	if len(m.Details.QueryID) == 0 {
		return errors.New("missing query key")
	}

	return nil
}

type ErrorMessage struct {
	Status  string      `json:"status"`
	Details interface{} `json:"details"`
}

func CreateSuccessMessage(details interface{}) (errMsg ErrorMessage) {
	return ErrorMessage{
		Status:  "success",
		Details: details,
	}
}

func CreateErrorMessage(details interface{}) (errMsg ErrorMessage) {
	return ErrorMessage{
		Status:  "error",
		Details: details,
	}
}
