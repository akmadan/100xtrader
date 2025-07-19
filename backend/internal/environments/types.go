package environments

type AgentConfig struct {
	Type           string `yaml:"type"`
	Count          int    `yaml:"count"`
	Aggressiveness string `yaml:"aggressiveness"`
}

type MarketEventConfig struct {
	Type      string `yaml:"type"`
	Frequency string `yaml:"frequency"`
}

type EnvironmentConfig struct {
	ID          string              `yaml:"id"`
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Volatility  string              `yaml:"volatility"`
	Trend       string              `yaml:"trend"`
	Liquidity   string              `yaml:"liquidity"`
	Agents      []AgentConfig       `yaml:"agents"`
	Events      []MarketEventConfig `yaml:"events"`
}
