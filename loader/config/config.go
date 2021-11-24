package config

// TODO: These should be injected to Starport's config.

// Config is glocal config values.
type Config struct {
	Plugins []Plugin `yaml:"plugins"`
}

// Plugin is a entry of plugin section.
type Plugin struct {
	Name string `yaml:"name"`
	Repo string `yaml:"repo"`
	Desc string `yaml:"desc"`
}
