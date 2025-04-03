package exchangerate

import (
	"currencyservice/internal/models"
	"currencyservice/internal/repo/currencies"
	"errors"
)

type Usecase struct {
	repo *currencies.Repo
}

func NewUsecase(repo *currencies.Repo) *Usecase {
	return &Usecase{repo: repo}
}

func (usecase Usecase) GetCurrency(code string) (models.Currency, error) {
	currency, err := usecase.repo.GetCurrencyByCode(code)
	if err != nil {
		return models.Currency{}, err
	}

	return currency, nil
}

func (usecase Usecase) GetAllCurrencies() ([]models.Currency, error) {
	currencies, err := usecase.repo.GetCurrencies()
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func (usecase Usecase) CreateNewCurrency(code, fullname, sign string) error {
	if err := usecase.repo.AddCurrency(code, fullname, sign); err != nil {
		return err
	}

	return nil
}

func (usecase Usecase) CreateExchangeRate(codeBaseCurrency, codeTargetCurrency string, rate float64) error {
	if _, err := usecase.repo.GetCurrencyByCode(codeBaseCurrency); err != nil {
		return err
	}

	if _, err := usecase.repo.GetCurrencyByCode(codeTargetCurrency); err != nil {
		return err
	}

	exchangerate, err := usecase.repo.GetExchangeRateByCodesPair(codeBaseCurrency, codeTargetCurrency)
	if err != nil && !errors.Is(err, models.ErrorExchangeRateNotFound) {
		return err
	}

	if exchangerate.ID != 0 {
		return models.ErrorExchangeRateAlreadyExists
	}

	if err := usecase.repo.AddExchangeRate(codeBaseCurrency, codeTargetCurrency, rate); err != nil {
		return err
	}

	return nil
}

func (usecase Usecase) GetExchangeRates() ([]models.CurrencyExchange, error) {
	exchangerates, err := usecase.repo.GetExchangeRates()
	if err != nil {
		return nil, err
	}

	return exchangerates, nil
}

func (usecase Usecase) GetExchangeRateByCodesPair(codeBaseCurrency, codeTargetCurrency string) (models.CurrencyExchange, error) {
	exchangerate, err := usecase.repo.GetExchangeRateByCodesPair(codeBaseCurrency, codeTargetCurrency)
	if err != nil {
		return models.CurrencyExchange{}, err
	}

	return exchangerate, nil
}

func (usecase Usecase) AddExchangeRate(codeBaseCurrency, codeTargetCurrency string, rate float64) error {
	if err := usecase.repo.AddExchangeRate(codeBaseCurrency, codeTargetCurrency, rate); err != nil {
		return err
	}

	return nil
}

func (usecase Usecase) UpdateExchangeRate(codeBaseCurrency, codeTargetCurrency string, newRate float64) error {
	if err := usecase.repo.UpdateExchangeRate(codeBaseCurrency, codeTargetCurrency, newRate); err != nil {
		return err
	}

	return nil
}

// GET /exchange?from=BASE_CURRENCY_CODE&to=TARGET_CURRENCY_CODE&amount=$AMOUNT
func (usecase Usecase) GetExchangeCurrencies(codeBaseCurrency, codeTargetCurrency string, amount float64) (models.GetExchangeCurrencies, error) {
	baseCurrency, err := usecase.repo.GetCurrencyByCode(codeBaseCurrency)
	if err != nil {
		return models.GetExchangeCurrencies{}, err
	}

	targetCurrency, err := usecase.repo.GetCurrencyByCode(codeTargetCurrency)
	if err != nil {
		return models.GetExchangeCurrencies{}, err
	}

	exchangerate, err := usecase.GetExchangeRateByCodesPair(codeBaseCurrency, codeTargetCurrency)
	if err != nil {
		return models.GetExchangeCurrencies{}, err
	}

	return models.GetExchangeCurrencies{
		BaseCurrency:    baseCurrency,
		TargetCurrency:  targetCurrency,
		Rate:            exchangerate.Rate,
		Amount:          amount,
		ConvertedAmount: amount * exchangerate.Rate,
	}, nil
}
