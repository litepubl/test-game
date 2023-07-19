package logger

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configurator  is a logger builder.
type Configurator struct {
	debugFile *os.File
	errorFile *os.File
}

// NewConfigurator is a logger constructor.
func NewConfigurator(config *Config) (*Configurator, error) {
	var (
		c   Configurator
		err error
	)

	c.errorFile, err = c.openFile(config.Path + config.ErrorFile)
	if err != nil {
		return nil, fmt.Errorf("NewConfigurator open errorFile %w", err)
	}

	c.debugFile, err = c.openFile(config.Path + config.DebugFile)
	if err != nil {
		return nil, fmt.Errorf("NewConfigurator open debugFile %w", err)
	}

	// consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	errorLevelWriter := &LevelWriter{c.errorFile, zerolog.ErrorLevel}
	debugLevelWriter := &LevelWriter{c.debugFile, zerolog.DebugLevel}

	multi := zerolog.MultiLevelWriter(os.Stdout, errorLevelWriter, debugLevelWriter)
	skipFrameCount := 2
	log.Logger = zerolog.New(multi).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	log.Info().Msg("Logger initialized")

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	return &c, nil
}

// Close the opened log files.
func (c *Configurator) Close() {
	log.Info().Msg("Logger closing")
	log.Logger = zerolog.Nop().With().Logger()

	if c.errorFile != nil {
		_ = c.errorFile.Close()
		c.errorFile = nil
	}

	if c.debugFile != nil {
		_ = c.debugFile.Close()
		c.debugFile = nil
	}
}

func (*Configurator) openFile(fileName string) (*os.File, error) {
	return os.OpenFile(
		fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o664,
	)
}
