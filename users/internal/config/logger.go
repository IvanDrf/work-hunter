package config

// Logger configuration
type LoggerConfig struct {
	Level      string   `yaml:"level"`
	Format     string   `yaml:"format"`
	OutputPath []string `yaml:"output_path"`
}
