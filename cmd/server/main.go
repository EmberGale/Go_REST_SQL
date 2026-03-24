package main

import (
	"GoRestSQL/internal/handler"
	"GoRestSQL/internal/repository"
	"GoRestSQL/internal/service"
	"database/sql"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

const (
	schemaSQL = `CREATE TABLE IF NOT EXISTS payments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	person TEXT,
	amount REAL,
	currency TEXT,
	time INTEGER,
	);`

	dbfile = "payments.db"
)

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err = db.Exec(schemaSQL); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := initDB(dbfile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	paymentRepo := repository.NewSqlitePaymentRepository(db)
	paymentService := service.NewPaymentServiceImpl(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	router := handler.NewRouter(paymentHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
