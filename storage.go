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
	connSter := "user=policies dbname=postgres password=aquaAssignment sslmode=disable"
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
