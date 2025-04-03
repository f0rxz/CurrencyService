package currencies

import (
	"currencyservice/internal/models"
	"database/sql"
	"errors"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

// POST /currencies
func (repo *Repo) AddCurrency(code, fullname, sign string) error {
	query := `
		INSERT INTO Currencies (Code, FullName, Sign)
		VALUES (?, ?, ?)
	`
	if _, err := repo.db.Exec(query, code, fullname, sign); err != nil {
		return err
	}

	return nil
}

// GET /currency/EUR
func (repo *Repo) GetCurrencyByID(id int) (models.Currency, error) {
	query := `
		SELECT ID, Code, FullName, Sign FROM Currencies
		WHERE ID=?
	`
	result, err := repo.db.Query(query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Currency{}, models.ErrorCurrencyNotFound
		}
		return models.Currency{}, err
	}

	currency := models.Currency{}
	if err := result.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign); err != nil {
		return models.Currency{}, err
	}

	return currency, nil
}

func (repo *Repo) GetCurrencyByCode(code string) (models.Currency, error) {
	query := `
		SELECT ID, Code, FullName, Sign FROM Currencies
		WHERE Code=?
	`
	result, err := repo.db.Query(query, code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Currency{}, models.ErrorCurrencyNotFound
		}
		return models.Currency{}, err
	}

	currency := models.Currency{}
	if err := result.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign); err != nil {
		return models.Currency{}, err
	}

	return currency, nil
}

// GET /currencies
func (repo *Repo) GetCurrencies() ([]models.Currency, error) {
	query := `
		SELECT ID, Code, FullName, Sign FROM Currencies
	`

	result, err := repo.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrorCurrencyNotFound
		}
		return nil, err
	}

	currencies := make([]models.Currency, 0, 10)
	defer result.Close()

	for result.Next() {
		currency := models.Currency{}
		if err := result.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}

// GET /exchangeRates
func (repo *Repo) GetExchangeRates() ([]models.CurrencyExchange, error) {
	query := `
		SELECT ID, BaseCurrencyId, TargetCurrencyId, Rate FROM ExchangeRates
	`

	result, err := repo.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrorCurrencyNotFound
		}
		return nil, err
	}

	exchangerates := make([]models.CurrencyExchange, 0, 10)
	defer result.Close()

	for result.Next() {
		exchangerate := models.CurrencyExchange{}
		if err := result.Scan(&exchangerate.ID, &exchangerate.BaseCurrencyId, &exchangerate.TargetCurrencyId, &exchangerate.Rate); err != nil {
			return nil, err
		}
		exchangerates = append(exchangerates, exchangerate)
	}

	return exchangerates, nil
}

// GET /exchangeRate/USDRUB
func (repo *Repo) GetExchangeRateByID(idBaseCurrency, idTargetCurrency int) (models.CurrencyExchange, error) {
	query := `
		SELECT ID, BaseCurrencyId, TargetCurrencyId, Rate FROM ExchangeRates
		WHERE BaseCurrencyId=? AND TargetCurrencyId=?
	`
	result, err := repo.db.Query(query, idBaseCurrency, idTargetCurrency)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CurrencyExchange{}, models.ErrorCurrencyNotFound
		}
		return models.CurrencyExchange{}, err
	}

	exchangerate := models.CurrencyExchange{}
	if err := result.Scan(&exchangerate.ID, &exchangerate.BaseCurrencyId, &exchangerate.TargetCurrencyId, &exchangerate.Rate); err != nil {
		return models.CurrencyExchange{}, err
	}

	return exchangerate, nil
}

// POST /exchangeRates
func (repo *Repo) AddExchangeRate(idBaseCurrency, idTargetCurrency int, rate float64) error {
	if _, err := repo.GetCurrencyByID(idBaseCurrency); err != nil {
		return err
	}

	if _, err := repo.GetCurrencyByID(idTargetCurrency); err != nil {
		return err
	}

	query := `
		INSERT INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
		VALUES (?, ?, ?) 
	`
	if _, err := repo.db.Exec(query, idBaseCurrency, idTargetCurrency, rate); err != nil {
		return err
	}

	return nil
}

// PATCH /exchangeRate/USDRUB
func (repo *Repo) UpdateExchangeRate(idBaseCurrency, idTargetCurrency int, newRate float64) error {
	query := `
		UPDATE ExchangeRates SET Rate = ? WHERE BaseCurrencyId = ? AND TargetCurrencyId = ?
	`

	if _, err := repo.db.Exec(query, newRate, idBaseCurrency, idTargetCurrency); err != nil {
		return err
	}

	return nil
}
