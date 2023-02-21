package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// 1. [] предполагается, что база данных и таблица существуют, и код потерпит неудачу, если это не так;
// 2. [] строка подключения жестко запрограммирована, даже пароль;
// 3. [] по-прежнему отсутствует метод Close, который закрывал бы соединение;
// 4. [] служба может закрыться, когда события все еще находятся в буфере записи: в этом случае события могут быть потеряны;
// 5. [] журнал хранит записи об удаленных значениях: он будет неограниченно расти.

type PostgresDBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type PostgresTransactionLogger struct {
	events chan<- Event
	errors <-chan error
	db     *sql.DB
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions (event_type, key, value) VALUES ($1, $2, $3)`

		for e := range events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()

		e := Event{}

		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error) {
	conStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		config.host, config.dbName, config.user, config.password)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}

	if err = logger.createTableIfNotExist(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) createTableIfNotExist() error {
	var err error

	query := `CREATE TABLE IF NOT EXISTS transactions (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	);`

	_, err = l.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
