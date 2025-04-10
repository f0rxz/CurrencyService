package handlers

import (
	"currencyservice/internal/controller/httpservice/handlers/exchanges"
	"currencyservice/internal/usecase/exchangerate"
)

type Handlers struct {
	ExchangesHandler *exchanges.Handler
}

func New(exchangeUsecase *exchangerate.Usecase) *Handlers {
	return &Handlers{
		ExchangesHandler: exchanges.NewHandler(exchangeUsecase),
	}
}
