package database

import (
"database/sql"
"fmt"
"log"
"os"

_ "github.com/lib/pq"
)

// InitDB - инициализирует подключение к PostgreSQL
func InitDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/todo_auth?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")

	// Создаем таблицы если их нет
	createTables(db)

	return db
}

// createTables - создает таблицы в БД
func createTables(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
id SERIAL PRIMARY KEY,
title VARCHAR(255) NOT NULL,
description TEXT,
completed BOOLEAN DEFAULT false,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

	CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
email VARCHAR(255) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
	`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	fmt.Println("Tables created successfully")
}
