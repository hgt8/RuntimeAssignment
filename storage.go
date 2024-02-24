package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type Storage interface {
	CreatePolicy(policy *CreatePolicyRequest) error
	UpdatePolicy(id int, policy *UpdatePolicyRequest) error //maybe by id and use patch?
	DeletePolicy(int) error
	GetPolicy(int) (*Policy, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func PostgresStore() (*PostgresStorage, error) {
	driverName := os.Getenv("PostgresDriverName")
	connStr := os.Getenv("PostgresConnectionString")
	db, err := sql.Open(driverName, connStr)
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

func (s *PostgresStorage) InitializeStorage() error {
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

func (s *PostgresStorage) UpdatePolicy(id int, policy *UpdatePolicyRequest) error {
	//goland:noinspection SqlNoDataSourceInspection
	sqlStatement := `
        UPDATE policies
        SET name = $1, author = $2, controls = $3, updated_at = $4
        WHERE id = $5;`
	_, err := s.db.Exec(sqlStatement, policy.PolicyName, policy.Author, policy.ControlData, time.Now().Add(2*time.Hour), id)
	return err
}

func (s *PostgresStorage) DeletePolicy(id int) error {
	//goland:noinspection SqlNoDataSourceInspection
	sqlStatement := `DELETE FROM policies WHERE id = $1;`
	result, err := s.db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (s *PostgresStorage) GetPolicy(id int) (*Policy, error) {
	//goland:noinspection SqlNoDataSourceInspection
	var policy Policy
	sqlStatement := `SELECT id, name, author, controls FROM policies WHERE id=$1;`
	row := s.db.QueryRow(sqlStatement, id)
	err := row.Scan(&policy.ID, &policy.PolicyName, &policy.Author, &policy.ControlData)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}
