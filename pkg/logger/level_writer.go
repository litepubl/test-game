package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type LevelWriter struct {
	io.Writer
	Level zerolog.Level
}

var _ zerolog.LevelWriter = &LevelWriter{}

func (lw *LevelWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l >= lw.Level { // Notice that it's ">=", not ">"
		return lw.Writer.Write(p)
	}

	return len(p), nil
}

func (lw *LevelWriter) Close() error {
	if c, ok := lw.Writer.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
