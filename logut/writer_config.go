package logut

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gitchander/logl"
)

type LogConfig struct {
	LogLevel     logl.Level   `json:"log-level"     yaml:"log-level"     toml:"log-level"`
	WriterConfig WriterConfig `json:"writer-config" yaml:"writer-config" toml:"writer-config"`
}

type WriterConfig struct {
	EnableStdout bool       `json:"enable-stdout" yaml:"enable-stdout" toml:"enable-stdout"`
	EnableFile   bool       `json:"enable-file"   yaml:"enable-file"   toml:"enable-file"`
	FileConfig   FileConfig `json:"file-config"   yaml:"file-config"   toml:"file-config"`
}

// lumberjack.Logger
type FileConfig struct {
	Filename   string `json:"filename"    yaml:"filename"    toml:"filename"`
	MaxSize    int    `json:"max-size"    yaml:"max-size"    toml:"max-size"`
	MaxAge     int    `json:"max-age"     yaml:"max-age"     toml:"max-age"`
	MaxBackups int    `json:"max-backups" yaml:"max-backups" toml:"max-backups"`
	LocalTime  bool   `json:"localtime"   yaml:"localtime"   toml:"localtime"`
	Compress   bool   `json:"compress"    yaml:"compress"    toml:"compress"`
}

var (
	DefaultFileConfig = FileConfig{
		Filename:   "logs/test.log",
		MaxSize:    50, // Maximum size in megabytes
		MaxAge:     30, // days
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   true,
	}

	DefaultLogConfig = LogConfig{
		LogLevel: logl.LevelDebug,
		WriterConfig: WriterConfig{
			EnableStdout: true,
			EnableFile:   true,
			FileConfig:   DefaultFileConfig,
		},
	}
)

// ------------------------------------------------------------------------------
func newLumberjack(fc FileConfig) LogWriter {
	return &lumberjack.Logger{
		Filename:   fc.Filename,
		MaxSize:    fc.MaxSize,
		MaxAge:     fc.MaxAge,
		MaxBackups: fc.MaxBackups,
		LocalTime:  fc.LocalTime,
		Compress:   fc.Compress,
	}
}

func _() {
	l := newLumberjack(DefaultFileConfig)

	var lw LogWriter = l
	defer lw.Close()

	log.SetOutput(lw)
	log.Println("Hello, lumberjack!")

	lw.Rotate()
}
