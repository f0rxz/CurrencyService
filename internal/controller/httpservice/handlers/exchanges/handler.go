package exchanges

import (
	"currencyservice/internal/usecase/exchangerate"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	exchangeUsecase *exchangerate.Usecase
}

func NewHandler(exchangeUsecase *exchangerate.Usecase) *Handler {
	return &Handler{
		exchangeUsecase: exchangeUsecase,
	}
}

func (h Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currencies, err := h.exchangeUsecase.GetAllCurrencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}

func (h Handler) GetCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Path[len("/currency/"):]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	currency, err := h.exchangeUsecase.GetCurrency(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currency)
}

func (h Handler) CreateNewCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	fullName := r.FormValue("fullname")
	sign := r.FormValue("sign")

	if code == "" || fullName == "" || sign == "" {
		http.Error(w, "All fields (code, fullname, sign) are required", http.StatusBadRequest)
		return
	}

	if err := h.exchangeUsecase.CreateNewCurrency(code, fullName, sign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Currency created successfully",
		"code":    code,
	})
}

func (h Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rates, err := h.exchangeUsecase.GetExchangeRates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, 0, len(rates))
	for _, rate := range rates {
		fmt.Println(rate)
		baseCurrency, err := h.exchangeUsecase.GetCurrency(rate.BaseCurrencyCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		targetCurrency, err := h.exchangeUsecase.GetCurrency(rate.TargetCurrencyCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response = append(response, map[string]interface{}{
			"id": rate.ID,
			"baseCurrency": map[string]interface{}{
				"id":   baseCurrency.ID,
				"name": baseCurrency.FullName,
				"code": baseCurrency.Code,
				"sign": baseCurrency.Sign,
			},
			"targetCurrency": map[string]interface{}{
				"id":   targetCurrency.ID,
				"name": targetCurrency.FullName,
				"code": targetCurrency.Code,
				"sign": targetCurrency.Sign,
			},
			"rate": rate.Rate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h Handler) GetExchangeRateByCodesPair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Path[len("/exchangeRate/"):]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	base := code[:3]
	target := code[3:]
	if base == "" || target == "" {
		http.Error(w, "Both base and target currency codes are required", http.StatusBadRequest)
		return
	}

	rate, err := h.exchangeUsecase.GetExchangeRateByCodesPair(base, target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}

func (h Handler) CreateExchangeRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	base := r.FormValue("base")
	target := r.FormValue("target")
	rate := r.FormValue("rate")

	if base == "" || target == "" || rate == "" {
		http.Error(w, "All fields (base, target, rate) are required", http.StatusBadRequest)
		return
	}

	rateValue, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		http.Error(w, "Invalid rate format", http.StatusBadRequest)
		return
	}

	if err := h.exchangeUsecase.CreateExchangeRate(base, target, rateValue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Exchange rate created successfully",
		"pair":    base + "/" + target,
	})
}

func (h Handler) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code := r.URL.Path[len("/exchangeRate/"):]
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	base := code[:3]
	target := code[3:]

	if base == "" || target == "" {
		http.Error(w, "Both base and target currency codes are required", http.StatusBadRequest)
		return
	}

	newRate := r.FormValue("newRate")

	if newRate == "" {
		http.Error(w, "Аield newRate are required", http.StatusBadRequest)
		return
	}

	rateValue, err := strconv.ParseFloat(newRate, 64)
	if err != nil {
		http.Error(w, "Invalid rate format", http.StatusBadRequest)
		return
	}

	if err := h.exchangeUsecase.UpdateExchangeRate(base, target, rateValue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Exchange rate updated successfully",
		"pair":    base + "/" + target,
	})
}

func (h Handler) GetExchangeCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amount")

	if from == "" || to == "" || amount == "" {
		http.Error(w, "All fields (from, to, amount) are required", http.StatusBadRequest)
		return
	}

	amountValue, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		http.Error(w, "Invalid rate format", http.StatusBadRequest)
		return
	}

	result, err := h.exchangeUsecase.GetExchangeCurrencies(from, to, amountValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
