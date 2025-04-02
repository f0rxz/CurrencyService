package models

import "errors"

type Currency struct {
	ID       int
	Code     string
	FullName string
	Sign     string
}

var (
	ErrorCurrencyNotFound = errors.New("Currency Not Found")
)
