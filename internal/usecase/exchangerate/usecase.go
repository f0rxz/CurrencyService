package exchangerate

import (
	"currencyservice/internal/models"
	"currencyservice/internal/repo/currencies"
)

type Usecase struct {
	repo *currencies.Repo
}

func NewUsecase(repo *currencies.Repo) *Usecase {
	return &Usecase{repo: repo}
}

func (usecase Usecase) GetCurrency(id int) (models.Currency, error) {
	currency, err := usecase.repo.GetCurrencyByID(id)

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

func (usecase Usecase) CreateExchangeRate(idBaseCurrency int, idTargetCurrency int, rate float64) error {
	if _, err := usecase.repo.GetCurrencyByID(idBaseCurrency); err != nil {
		return err
	}

	if _, err := usecase.repo.GetCurrencyByID(idTargetCurrency); err != nil {
		return err
	}

	if _, err := usecase.repo.GetExchangeRateByID(idBaseCurrency, idTargetCurrency); err != nil {
		if err := usecase.repo.AddExchangeRate(idBaseCurrency, idTargetCurrency, rate); err != nil {
			return err
		}
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

func (usecase Usecase) GetExchangeRateByID(idBaseCurrency, idTargetCurrency int) (models.CurrencyExchange, error) {
	exchangerate, err := usecase.repo.GetExchangeRateByID(idBaseCurrency, idTargetCurrency)

	if err != nil {
		return models.CurrencyExchange{}, err
	}

	return exchangerate, nil
}

func (usecase Usecase) AddExchangeRate(idBaseCurrency, idTargetCurrency int, rate float64) error {
	if err := usecase.repo.AddExchangeRate(idBaseCurrency, idTargetCurrency, rate); err != nil {
		return err
	}

	return nil
}

func (usecase Usecase) UpdateExchangeRate(idBaseCurrency, idTargetCurrency int, newRate float64) error {
	if err := usecase.repo.UpdateExchangeRate(idBaseCurrency, idTargetCurrency, newRate); err != nil {
		return err
	}
	return nil
}
