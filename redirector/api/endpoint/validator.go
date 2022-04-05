package endpoint

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func ValidateSavingsPeriod(period string) (err error) {
	if val, err := strconv.ParseInt(period, 10, 64); err != nil {
		return errors.New("invalid savings period")
	} else if val <= 0 {
		return errors.New("invalid savings period")
	}
	return nil
}

func ValidateAmount(amount string) (err error) {
	if val, err := strconv.ParseFloat(amount, 64); err != nil {
		return errors.New("invalid amount")
	} else if val <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}

func ValidateCurrency(currency string) (err error) {
	currencyU := strings.ToUpper(currency)
	if strings.Compare(currencyU, "VND") != 0 &&
		strings.Compare(currencyU, "USD") != 0 {
		// currently only support these 2 currency
		return errors.New("invalid currency unit")
	}
	return nil
}

func ValidateTime(t string) (err error) {
	_, err = time.Parse(t, "Mon, 02 Jan 2006 15:04:05 MST")
	if err != nil {
		return errors.New("invalid time format")
	}
	return nil
}

func ValidateInterestRate(rate string) (err error) {
	if val, err := strconv.ParseFloat(rate, 64); err != nil {
		return errors.New("invalid interest rate")
	} else if val <= 0 {
		return errors.New("invalid interest rate")
	}
	return nil
}
