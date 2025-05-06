package main

import (
	"database/sql"
	"embed"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// openDB открывает базу данных SQLite по пути path и возвращает *sql.DB
// или вызывает log.Fatal в случае ошибки.
func openDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	// Проверим подключение
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("✅ Подключение к базе данных установлено")
	return db
}

// applyMigrations применяет миграции к базе данных, используя миграции,
// хранящиеся в migrationFiles.
func applyMigrations(db *sql.DB) error {
	log.Println("🔄 Применяем миграции к базе данных…")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("✅ Миграции успешно применены")
	return nil
}

// durationSince возвращает длительность времени, прошедшего с момента времени t.
// Если t.Valid == false, то возвращается nil.
// Результат округляется до ближайшей секунды.
func durationSince(t sql.NullTime) *time.Duration {
	if t.Valid {
		d := time.Since(t.Time).Round(time.Second)
		return &d
	}
	return nil
}

// isSettingEnabled возвращает true, если значение настройки с указанным key
// равно "true", иначе false. Если настройки с указанным key не существует,
// то возвращается false, nil. Если произошла какая-либо ошибка, то
// возвращается false, error.
func isSettingEnabled(db *sql.DB, key string) (bool, error) {
	var value string
	err := db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return value == "true", nil
}

// flipSetting инвертирует булево значение настройки с указанным key
// в базе данных. Если настройки с указанным key не существует, то
// ничего не происходит. Если произошла какая-либо ошибка, то
// возвращается error.
func flipSetting(db *sql.DB, key string) error {
	_, err := db.Exec(`UPDATE settings SET value = CASE WHEN value = 'true' THEN 'false' ELSE 'true' END WHERE key = ?`, key)
	return err
}

// getSetting возвращает значение настройки с указанным key.
// Если настройки с указанным key не существует, то возвращается пустая строка, nil.
// Если произошла какая-либо ошибка, то возвращается пустая строка, error.
func getSetting(db *sql.DB, key string) (string, error) {
	var value string
	err := db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	return value, err
}
