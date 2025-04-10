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
// func (repo *Repo) GetCurrencyByID(id int) (models.Currency, error) {
// 	query := `
// 		SELECT ID, Code, FullName, Sign FROM Currencies
// 		WHERE ID=?
// 	`
// 	result, err := repo.db.Query(query, id)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return models.Currency{}, models.ErrorCurrencyNotFound
// 		}
// 		return models.Currency{}, err
// 	}

// 	currency := models.Currency{}
// 	if err := result.Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign); err != nil {
// 		return models.Currency{}, err
// 	}

// 	return currency, nil
// }

func (repo *Repo) GetCurrencyByCode(code string) (models.Currency, error) {
	query := `
        SELECT ID, Code, FullName, Sign FROM Currencies
        WHERE Code=?
    `
	currency := models.Currency{}
	err := repo.db.QueryRow(query, code).Scan(&currency.ID, &currency.Code, &currency.FullName, &currency.Sign)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Currency{}, models.ErrorCurrencyNotFound
		}
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
		SELECT ID, BaseCurrencyCode, TargetCurrencyCode, Rate FROM ExchangeRates
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
		if err := result.Scan(&exchangerate.ID, &exchangerate.BaseCurrencyCode, &exchangerate.TargetCurrencyCode, &exchangerate.Rate); err != nil {
			return nil, err
		}
		exchangerates = append(exchangerates, exchangerate)
	}

	return exchangerates, nil
}

// GET /exchangeRate/USDRUB
func (repo *Repo) GetExchangeRateByCodesPair(codeBaseCurrency, codeTargetCurrency string) (models.CurrencyExchange, error) {
	baseCurrency, err := repo.GetCurrencyByCode(codeBaseCurrency)
	if err != nil {
		return models.CurrencyExchange{}, err
	}

	targetCurrency, err := repo.GetCurrencyByCode(codeTargetCurrency)
	if err != nil {
		return models.CurrencyExchange{}, err
	}

	query := `
        SELECT ID, BaseCurrencyCode, TargetCurrencyCode, Rate FROM ExchangeRates
        WHERE BaseCurrencyCode=? AND TargetCurrencyCode=?
    `

	exchangerate := models.CurrencyExchange{}
	err = repo.db.QueryRow(query, baseCurrency.Code, targetCurrency.Code).Scan(
		&exchangerate.ID,
		&exchangerate.BaseCurrencyCode,
		&exchangerate.TargetCurrencyCode,
		&exchangerate.Rate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CurrencyExchange{}, models.ErrorExchangeRateNotFound
		}
		return models.CurrencyExchange{}, err
	}

	return exchangerate, nil
}

// POST /exchangeRates
func (repo *Repo) AddExchangeRate(codeBaseCurrency, codeTargetCurrency string, rate float64) error {
	baseCurrency, err := repo.GetCurrencyByCode(codeBaseCurrency)
	if err != nil {
		return err
	}

	targetCurrency, err := repo.GetCurrencyByCode(codeTargetCurrency)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO ExchangeRates (BaseCurrencyCode, TargetCurrencyCode, Rate)
		VALUES (?, ?, ?) 
	`
	if _, err := repo.db.Exec(query, baseCurrency.Code, targetCurrency.Code, rate); err != nil {
		return err
	}

	return nil
}

// PATCH /exchangeRate/USDRUB
func (repo *Repo) UpdateExchangeRate(codeBaseCurrency, codeTargetCurrency string, newRate float64) error {
	baseCurrency, err := repo.GetCurrencyByCode(codeBaseCurrency)
	if err != nil {
		return err
	}

	targetCurrency, err := repo.GetCurrencyByCode(codeTargetCurrency)
	if err != nil {
		return err
	}

	query := `
		UPDATE ExchangeRates SET Rate = ? WHERE BaseCurrencyCode = ? AND TargetCurrencyCode = ?
	`

	if _, err := repo.db.Exec(query, newRate, baseCurrency.Code, targetCurrency.Code); err != nil {
		return err
	}

	return nil
}
