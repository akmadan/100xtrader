package engine

import (
	"sync"

	"github.com/akshitmadan/100xtrader/backend/internal/data"
	"github.com/akshitmadan/100xtrader/backend/internal/data/repos"
)

type EnvironmentManager struct {
	mu                 sync.RWMutex
	environments       map[string]*data.Environment
	currentEnvironment string
	repo               repos.EnvironmentRepository
}

var (
	envManager     *EnvironmentManager
	envManagerOnce sync.Once
)

// GetEnvironmentManager returns the singleton instance
func GetEnvironmentManager(repo repos.EnvironmentRepository) *EnvironmentManager {
	envManagerOnce.Do(func() {
		envManager = &EnvironmentManager{
			environments: make(map[string]*data.Environment),
			repo:         repo,
		}
		envManager.LoadEnvironmentsFromRepo()
		envManager.LoadCurrentEnvironmentFromRepo()
	})
	return envManager
}

// AddEnvironment adds a new environment and persists it
func (m *EnvironmentManager) AddEnvironment(env *data.Environment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.environments[env.ID] = env
	return m.repo.SaveEnvironment(env)
}

// ListEnvironments returns all environments
func (m *EnvironmentManager) ListEnvironments() []*data.Environment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var envs []*data.Environment
	for _, env := range m.environments {
		envs = append(envs, env)
	}
	return envs
}

// SetCurrentEnvironment sets and persists the current environment
func (m *EnvironmentManager) SetCurrentEnvironment(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentEnvironment = id
	return m.repo.SaveCurrentEnvironment(id)
}

// GetCurrentEnvironment returns the current environment
func (m *EnvironmentManager) GetCurrentEnvironment() *data.Environment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.environments[m.currentEnvironment]
}

// LoadEnvironmentsFromRepo loads environments from the repository
func (m *EnvironmentManager) LoadEnvironmentsFromRepo() error {
	envs, err := m.repo.ListEnvironments()
	if err != nil {
		return err
	}
	for _, env := range envs {
		m.environments[env.ID] = env
	}
	return nil
}

// LoadCurrentEnvironmentFromRepo loads the current environment selection from the repository
func (m *EnvironmentManager) LoadCurrentEnvironmentFromRepo() error {
	id, err := m.repo.GetCurrentEnvironment()
	if err != nil {
		return err
	}
	m.currentEnvironment = id
	return nil
}
