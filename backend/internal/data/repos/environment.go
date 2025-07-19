package repos

import (
	"database/sql"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
)

type EnvironmentRepository interface {
	SaveEnvironment(env *data.Environment) error
	ListEnvironments() ([]*data.Environment, error)
	GetEnvironmentByID(id string) (*data.Environment, error)
	SaveCurrentEnvironment(id string) error
	GetCurrentEnvironment() (string, error)
}

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) SaveEnvironment(env *data.Environment) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO environments (id, name, description, volatility, trend, liquidity) VALUES (?, ?, ?, ?, ?, ?)`, env.ID, env.Name, env.Description, env.Volatility, env.Trend, env.Liquidity)
	return err
}

func (r *environmentRepository) ListEnvironments() ([]*data.Environment, error) {
	rows, err := r.db.Query(`SELECT id, name, description, volatility, trend, liquidity FROM environments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var envs []*data.Environment
	for rows.Next() {
		var env data.Environment
		if err := rows.Scan(&env.ID, &env.Name, &env.Description, &env.Volatility, &env.Trend, &env.Liquidity); err != nil {
			return nil, err
		}
		envs = append(envs, &env)
	}
	return envs, nil
}

func (r *environmentRepository) GetEnvironmentByID(id string) (*data.Environment, error) {
	row := r.db.QueryRow(`SELECT id, name, description, volatility, trend, liquidity FROM environments WHERE id = ?`, id)
	var env data.Environment
	if err := row.Scan(&env.ID, &env.Name, &env.Description, &env.Volatility, &env.Trend, &env.Liquidity); err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *environmentRepository) SaveCurrentEnvironment(id string) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES ('current_environment', ?)`, id)
	return err
}

func (r *environmentRepository) GetCurrentEnvironment() (string, error) {
	row := r.db.QueryRow(`SELECT value FROM meta WHERE key = 'current_environment'`)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
