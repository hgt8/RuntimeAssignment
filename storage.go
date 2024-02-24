package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	CreatePolicy(policy *CreatePolicyRequest) error
	UpdatePolicy(policy *Policy) error //maybe by id and use patch?
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
	//goland:noinspection SqlNoDataSourceInspection
	sqlStatement := `
		CREATE TABLE IF NOT EXISTS policies
		(id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		author VARCHAR(255) NOT NULL,
		controls JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)
		`
	_, err := s.db.Exec(sqlStatement)
	return err
}

func (s *PostgresStorage) CreatePolicy(policy *CreatePolicyRequest) error {
	//goland:noinspection SqlNoDataSourceInspection
	sqlStatement :=
		`INSERT INTO policies (name, author, controls)
		 VALUES ($1, $2, $3)`
	_, err := s.db.Exec(sqlStatement, policy.PolicyName, policy.Author, policy.ControlData)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *PostgresStorage) UpdatePolicy(policy *Policy) error {
	return nil
}

func (s *PostgresStorage) DeletePolicy(id int) error {
	return nil
}

func (s *PostgresStorage) GetPolicy(id int) (*Policy, error) {
	return nil, nil
}
