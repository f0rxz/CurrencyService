package models

import "errors"

type CurrencyExchange struct {
	ID               int
	BaseCurrencyId   int
	TargetCurrencyId int
	Rate             float64
}

var (
	ErrorExchangeRateNotFound = errors.New("Exchange rate Not Found")
)
