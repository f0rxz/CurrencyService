package repo

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Repo struct {
	db *sql.DB
}

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "sqlite-database.db")

	if err != nil {
		return nil, errors.New("Error connecting to database")
	}

	currencies := `CREATE TABLE IF NOT EXISTS Currencies (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Code VARCHAR(10) NOT NULL,
    FullName VARCHAR(100) NOT NULL,
    Sign VARCHAR(10) NOT NULL
	);`

	if _, err := db.Exec(currencies); err != nil {
		return nil, err
	}

	currencyexchange := `CREATE TABLE IF NOT EXISTS ExchangeRates (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    BaseCurrencyId INTEGER NOT NULL,
    TargetCurrencyId INTEGER NOT NULL,
    Rate DECIMAL(6, 4) NOT NULL,
    FOREIGN KEY (BaseCurrencyId) REFERENCES Currencies(ID) ON DELETE CASCADE,
    FOREIGN KEY (TargetCurrencyId) REFERENCES Currencies(ID) ON DELETE CASCADE
	);`

	if _, err := db.Exec(currencyexchange); err != nil {
		return nil, err
	}

	return db, nil
}

func (repo *Repo) Close() {
	repo.db.Close()
}
