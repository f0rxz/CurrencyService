package models

import (
	"errors"
)

var (
	ErrorExchangeRateNotFound      = errors.New("Exchange rate Not Found")
	ErrorExchangeRateAlreadyExists = errors.New("Exchange rate already exists")
)

type CurrencyExchange struct {
	ID               int
	BaseCurrencyId   int
	TargetCurrencyId int
	Rate             float64
}

type GetExchangeCurrencies struct {
	BaseCurrency    Currency
	TargetCurrency  Currency
	Rate            float64
	Amount          float64
	ConvertedAmount float64
}
