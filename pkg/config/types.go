package config

// RulesConfig contains all the rules
type RulesConfig struct {
	Rules []Rule `yaml:"rules"`
}

// Rule defines a rule
type Rule struct {
	Name      string   `yaml:"name"`
	Schedules []string `yaml:"schedules"`
	Repos     []string `yaml:"repos"`
	Channels  []string `yaml:"channels"`
	Ignore    Ignore   `yaml:"ignore,omitempty"`
}

// Ignore contains conditions on which PRs are ignored
type Ignore struct {
	YoungerThan string `yaml:"youngerThan,omitempty"`
	OlderThan   string `yaml:"olderThan,omitempty"`
}
