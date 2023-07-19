package logger

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

// LevelWriter  is simple wraper for log write by level.
type LevelWriter struct {
	io.Writer
	Level zerolog.Level
}

var _ zerolog.LevelWriter = (*LevelWriter)(nil)

// WriteLevel filter write by log level.
func (lw *LevelWriter) WriteLevel(l zerolog.Level, p []byte) (int, error) {
	if l >= lw.Level { // Notice that it's ">=", not ">"
		n, err := lw.Writer.Write(p)
		if err != nil {
			return n, fmt.Errorf("LevelWriter.WriteLevel  error %w", err)
		}

		return n, nil
	}

	return len(p), nil
}

// Close delegated writer.
func (lw *LevelWriter) Close() error {
	if c, ok := lw.Writer.(io.Closer); ok {
		err := c.Close()
		if err != nil {
			return fmt.Errorf("LevelWriter- Close error %w", err)
		}
	}

	return nil
}
