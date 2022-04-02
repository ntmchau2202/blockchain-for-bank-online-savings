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
	Details map[string]interface{} `json:"details"`
}

func (m AddSavingsAccountMessage) Validate() (err error) {

	if err = m.Message.Validate(); err != nil {
		return
	}

	if strings.Compare(m.Message.Command, AddSavingsAccountCmd.ToString()) != 0 {
		return errors.New("invalid command")
	}

	if len(m.Details["savingsaccount_number"].(string)) == 0 {
		return errors.New("missing saving account number")
	}

	if len(m.Details["owner_id"].(string)) == 0 {
		return errors.New("missing owner id")
	}

	if len(m.Details["owner_phone"].(string)) == 0 {
		return errors.New("missing customer phone for validation")
	}

	if len(m.Details["product_type"].(string)) == 0 {
		return errors.New("missing product type")
	}

	if _, err := strconv.ParseFloat(m.Details["savings_amount"].(string), 64); err != nil {
		return errors.New("invalid initial savings amount")
	}

	if _, err := strconv.ParseInt(m.Details["savings_period"].(string), 10, 64); err != nil {
		return errors.New("invalid saving period")
	}

	if _, err := strconv.ParseFloat(m.Details["interest_rate"].(string), 64); err != nil {
		return errors.New("invalid interest rate")
	}

	if _, err := strconv.ParseFloat(m.Details["estimated_interest_amount"].(string), 64); err != nil {
		return errors.New("invalid estimated insterest amount")
	}

	if err = validateTime(m.Details["open_time"].(string)); err != nil {
		return err
	}

	if len(m.Details["currency"].(string)) != 3 || strings.Compare(strings.ToUpper(m.Details["currency"].(string)), m.Details["currency"].(string)) != 0 {
		return errors.New("invalid transaction unit format")
	}
	return nil
}

type SettleSavingsAccountMessage struct {
	Message
	Details map[string]interface{} `json:"details"`
}

func (m SettleSavingsAccountMessage) Validate() (err error) {
	if err = m.Message.Validate(); err != nil {
		return
	}

	if strings.Compare(m.Message.Command, SettleSavingsAccountCmd.ToString()) != 0 {
		return errors.New("invalid command")
	}

	if len(m.Details["savingsaccount_id"].(string)) == 0 {
		return errors.New("missing savings account id")
	}

	if len(m.Details["owner_id"].(string)) == 0 {
		return errors.New("missing owner id")
	}

	if len(m.Details["owner_phone"].(string)) == 0 {
		return errors.New("missing owner phone for validation")
	}

	if len(m.Details["actual_interest_amount"].(string)) == 0 {
		return errors.New("missing actual interest amount")
	}

	if err = validateTime(m.Details["settle_time"].(string)); err != nil {
		return err
	}

	return nil
}

type QueryMessage struct {
	Message
	Details map[string]interface{} `json:"details"`
}

func (m QueryMessage) Validate() (err error) {
	if err = m.Message.Validate(); err != nil {
		return
	}

	if strings.Compare(m.Message.Command, SearchTxnsByBankAccountCmd.ToString()) != 0 && strings.Compare(m.Message.Command, SearchTxnsBySavingsAccountCmd.ToString()) != 0 {
		return errors.New("invalid command")
	}

	if len(m.Details["query_id"].(string)) == 0 {
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
