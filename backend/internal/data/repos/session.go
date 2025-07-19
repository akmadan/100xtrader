package repos

import (
	"database/sql"
	"time"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

type SessionRepository interface {
	CreateSession(session *data.Session) error
	EndSession(id string, endedAt time.Time) error
	ListSessions(user string) ([]*data.Session, error)
	GetSessionByID(id string) (*data.Session, error)
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(session *data.Session) error {
	_, err := r.db.Exec(`INSERT INTO sessions (id, user, environment, ticker, started_at) VALUES (?, ?, ?, ?, ?)`,
		session.ID, session.User, session.Environment, session.Ticker, session.StartedAt)
	return err
}

func (r *sessionRepository) EndSession(id string, endedAt time.Time) error {
	_, err := r.db.Exec(`UPDATE sessions SET ended_at = ? WHERE id = ?`, endedAt, id)
	return err
}

func (r *sessionRepository) ListSessions(user string) ([]*data.Session, error) {
	rows, err := r.db.Query(`SELECT id, user, environment, ticker, started_at, ended_at FROM sessions WHERE user = ? ORDER BY started_at DESC`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []*data.Session
	for rows.Next() {
		var s data.Session
		var endedAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.User, &s.Environment, &s.Ticker, &s.StartedAt, &endedAt); err != nil {
			return nil, err
		}
		if endedAt.Valid {
			s.EndedAt = endedAt.Time
		}
		sessions = append(sessions, &s)
	}
	return sessions, nil
}

func (r *sessionRepository) GetSessionByID(id string) (*data.Session, error) {
	row := r.db.QueryRow(`SELECT id, user, environment, ticker, started_at, ended_at FROM sessions WHERE id = ?`, id)
	var s data.Session
	var endedAt sql.NullTime
	if err := row.Scan(&s.ID, &s.User, &s.Environment, &s.Ticker, &s.StartedAt, &endedAt); err != nil {
		return nil, err
	}
	if endedAt.Valid {
		s.EndedAt = endedAt.Time
	}
	return &s, nil
}
