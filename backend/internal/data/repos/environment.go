package repos

import (
	"github.com/akshitmadan/100xtrader/backend/internal/environments"
)

type EnvironmentRepository interface {
	ListEnvironments() ([]environments.EnvironmentConfig, error)
	GetEnvironmentByID(id string) (*environments.EnvironmentConfig, error)
}

type environmentRepository struct {
	dir string
}

func NewEnvironmentRepository(dir string) EnvironmentRepository {
	return &environmentRepository{dir: dir}
}

func (r *environmentRepository) ListEnvironments() ([]environments.EnvironmentConfig, error) {
	return environments.LoadEnvironments(r.dir)
}

func (r *environmentRepository) GetEnvironmentByID(id string) (*environments.EnvironmentConfig, error) {
	envs, err := environments.LoadEnvironments(r.dir)
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if env.ID == id {
			return &env, nil
		}
	}
	return nil, nil
}
