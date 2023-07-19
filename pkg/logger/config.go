// Package logger contains cofiguration for logger and multi level writer.
package logger

// Config  contains arguments for initilizing logger.
type Config struct {
	Path      string `yaml:"path" env:"LOG_PATH" envDefault:"/logs/"`
	ErrorFile string `yaml:"errorFile" env:"LOG_ERROR_FILE" envDefault:"error.log"`
	DebugFile string `yaml:"debugFile" env:"LOG_DEBUG_FILE" envDefault:"debug.log"`
}
