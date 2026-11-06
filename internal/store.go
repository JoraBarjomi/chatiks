package internal

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Store struct {
	db *sql.DB
}

func InitStore(ctx context.Context, dsn string) (*Store, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.initSchema(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) initSchema(ctx context.Context) error {
	const q = `CREATE TABLE IF NOT EXISTS messages (
		id BIGSERIAL PRIMARY KEY,
		room TEXT NOT NULL,
		username TEXT NOT NULL,
		text TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);
	CREATE INDEX IF NOT EXISTS idx_messages_room ON messages (room, id);`
	_, err := s.db.ExecContext(ctx, q)
	return err
}

func (s *Store) SaveMessage(ctx context.Context, m Message) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO messages (room, username, text, created_at) VALUES ($1, $2, $3, $4)`, m.Room, m.User, m.Text, m.Time)
	return err
}

func (s *Store) LastMessages(ctx context.Context, room string, limit int) ([]Message, error) {

	rows, err := s.db.QueryContext(ctx, `SELECT room, username, text, created_at FROM 
	(SELECT id, room, username, text, created_at FROM messages WHERE room = $1 ORDER BY id DESC LIMIT $2) as recent ORDER BY id ASC`, room, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs := make([]Message, 0, limit)
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.Room, &m.User, &m.Text, &m.Time); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}
