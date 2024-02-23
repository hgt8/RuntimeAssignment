package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreatePolicy(policy *Policy) error
	UpdatePolicy(policy *Policy) error //maybe by id?
	DeletePolicy(int) error
	GetPolicy(int) (*Policy, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func PostgresStore() (*PostgresStorage, error) {
	connSter := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connSter)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	return s.createPoliciesTable()
}

func (s *PostgresStorage) createPoliciesTable() error {
	//backticks does not work for some reason, used
	sqlStatement := "CREATE TABLE IF NOT EXISTS policies \n " +
		"(id SERIAL PRIMARY KEY,\n    " +
		"name VARCHAR(255) UNIQUE NOT NULL,\n    " +
		"author VARCHAR(255) NOT NULL,\n    " +
		"controls JSONB NOT NULL,\n    " +
		"created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,\n" +
		"updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP\n)"

	_, err := s.db.Exec(sqlStatement)
	return err
}

func (s *PostgresStorage) CreatePolicy(*Policy) error {
	return nil
}

func (s *PostgresStorage) UpdatePolicy(*Policy) error {
	return nil
}

func (s *PostgresStorage) DeletePolicy(id int) error {
	return nil
}

func (s *PostgresStorage) GetPolicy(id int) (*Policy, error) {
	return nil, nil
}
