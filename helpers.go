package main

import (
	"database/sql"
	"embed"
	"log"

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
