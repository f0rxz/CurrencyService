package models

import "errors"

var (
	ErrorCurrencyNotFound = errors.New("Currency Not Found")
)

type Currency struct {
	ID       int
	Code     string
	FullName string
	Sign     string
}
