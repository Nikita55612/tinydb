package server

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbInit sync.Once
	dbConn *sql.DB
)

func InitDBConn() {
	dbInit.Do(func() {
		var err interface{}
		dbConn, err = newDBConn("./storage.db")
		if err != nil {
			panic(err)
		}
	})
}

func newDBConn(dbPath string) (*sql.DB, error) {
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	pragmas := []string{
		"PRAGMA journal_mode=WAL",        // Write-Ahead Logging для конкурентного доступа
		"PRAGMA synchronous=NORMAL",      // Баланс между скоростью и надежностью
		"PRAGMA cache_size=-1000000",     // Кэш примерно 1GB (-1000000 * 1KB)
		"PRAGMA page_size=4096",          // Оптимальный размер страницы
		"PRAGMA mmap_size=20000000000",   // Memory-mapped I/O ~20GB
		"PRAGMA wal_autocheckpoint=1000", // Автоматический чекпоинт WAL
		"PRAGMA busy_timeout=5000",       // Таймаут для блокировок
		"PRAGMA temp_store=MEMORY",       // Временные таблицы в памяти
		"PRAGMA foreign_keys=ON",         // Обеспечение целостности данных
	}

	for _, pragma := range pragmas {
		if _, err := dbConn.Exec(pragma); err != nil {
			return nil, err
		}
	}

	dbConn.SetMaxOpenConns(1)
	dbConn.SetMaxIdleConns(1)

	tx, err := dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS data (
            hash BLOB PRIMARY KEY,
            data BLOB NOT NULL
        )
    `)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	tx = nil

	return dbConn, nil
}

func Insert(hash []byte, data []byte) error {
	_, err := dbConn.Exec("INSERT INTO data (hash, data) VALUES (?, ?)", hash, data)

	return err
}

func Replace(hash []byte, new_hash []byte, data []byte) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	tx.Exec("DELETE FROM data WHERE hash = ?", hash)
	_, err = tx.Exec("INSERT INTO data (hash, data) VALUES (?, ?)", new_hash, data)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	tx = nil

	return err
}

func Select(hash []byte) ([]byte, error) {
	var data []byte
	err := dbConn.QueryRow("SELECT data FROM data WHERE hash = ?", hash).Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
