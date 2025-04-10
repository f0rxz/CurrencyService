package httpservice

import (
	"currencyservice/internal/controller/httpservice/handlers"
	"currencyservice/internal/usecase/exchangerate"
	"fmt"
	"net/http"
)

type Server struct {
	handlers *handlers.Handlers
}

func NewServer(exchangeUsecase *exchangerate.Usecase) *Server {
	return &Server{
		handlers: handlers.New(exchangeUsecase),
	}
}

func (s Server) Start(port uint16) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s Server) SetupRoutes() {
	http.HandleFunc("/currencies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handlers.ExchangesHandler.GetCurrencies(w, r)
		case http.MethodPost:
			s.handlers.ExchangesHandler.CreateNewCurrency(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/currency/", s.handlers.ExchangesHandler.GetCurrencyByCode)

	http.HandleFunc("/exchangeRates", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handlers.ExchangesHandler.GetExchangeRates(w, r)
		case http.MethodPost:
			s.handlers.ExchangesHandler.CreateExchangeRate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/exchangeRate/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handlers.ExchangesHandler.GetExchangeRateByCodesPair(w, r)
		case http.MethodPatch:
			s.handlers.ExchangesHandler.UpdateExchangeRate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/exchange", s.handlers.ExchangesHandler.GetExchangeCurrencies)

}
