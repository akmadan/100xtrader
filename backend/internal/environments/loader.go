package environments

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadEnvironments reads all YAML files in the given directory and returns a slice of EnvironmentConfig
func LoadEnvironments(dir string) ([]EnvironmentConfig, error) {
	var envs []EnvironmentConfig

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !(filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml") {
			continue
		}
		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var env EnvironmentConfig
		if err := yaml.Unmarshal(data, &env); err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}
	return envs, nil
}
