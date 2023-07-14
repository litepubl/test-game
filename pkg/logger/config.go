package logger

type Config struct {
	Path      string `yaml:"path" env:"LOG_PATH" envDefault:"/logs/"`
	ErrorFile string `yaml:"error_file" env:"LOG_ERROR_FILE" envDefault:"error.log"`
	DebugFile string `yaml:"debug_file" env:"LOG_DEBUG_FILE" envDefault:"debug.log"`
}
