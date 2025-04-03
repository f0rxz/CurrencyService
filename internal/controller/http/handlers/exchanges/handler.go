package exchanges

import (
	"currencyservice/internal/usecase/exchangerate"
	"net/http"
)

// here will be handlers
type Handler struct {
	exchangeUsecase *exchangerate.Usecase
}

func NewHandler(exchangeUsecase *exchangerate.Usecase) *Handler {
	return &Handler{
		exchangeUsecase: exchangeUsecase,
	}
}

func (h Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) GetCurrencyByCode(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) CreateNewCurrency(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) GetExchangeRateByCodesPair(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) CreateExchangeRate(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {

}
