package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var (
	Global zerolog.Logger

	logConfig = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Kitchen,
	}

	defaultOptions = &Options{
		MinimumLevel: zerolog.NoLevel,
		Callers:      true,
	}
)

type Options struct {
	MinimumLevel zerolog.Level
	Callers      bool
}

// init initializes the logger configuration.
func init() {
	logConfig.FormatLevel = func(i interface{}) string {
		switch i.(string) {
		case "info":
			return "\x1b[1m\x1b[38;2;95;250;213mINFO\x1b[0m \x1b[38;5;239m>\x1b[0m"
		case "debug":
			return "\x1b[1m\x1b[38;2;93;94;255mDEBU\x1b[0m \x1b[38;5;239m>\x1b[0m"
		case "warn":
			return "\x1b[1m\x1b[38;2;216;252;141mWARN\x1b[0m \x1b[38;5;239m>\x1b[0m"
		case "error":
			return "\x1b[1m\x1b[38;2;251;97;137mERRO\x1b[0m \x1b[38;5;239m>\x1b[0m"
		case "fatal":
			return "\x1b[1m\x1b[38;2;179;94;220mFATA\x1b[0m \x1b[38;5;239m>\x1b[0m"
		default:
			return i.(string)
		}
	}

	// Format of the caller thing aka Log.go:20 >
	logConfig.FormatCaller = func(i interface{}) string {
		return fmt.Sprintf("\x1b[0m%v \x1b[38;5;239m>\x1b[0m", i)
	}

	// File log
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file

		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}

		file = short
		return fmt.Sprintf("%v:%v", file, line)
	}

	// FormatFieldName is cyan like - we set it to a dark grayish tone.
	logConfig.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\x1b[38;5;239m%v=\x1b[0m", i)
	}

	// FormatErrFieldName is cyan like as well - we set it to a dark grayish tone.
	logConfig.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("\x1b[38;5;239m%v=\x1b[0m", i)
	}

	// FormatErrFieldValue is default red - we reset it.
	logConfig.FormatErrFieldValue = func(i interface{}) string {
		return i.(string)
	}

	logConfig.FormatMessage = func(i interface{}) string {
		if msg, ok := i.(string); ok {
			return msg
		}
		return ""
	}

	log.Logger = log.Output(logConfig).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.GlobalLevel()) /* zerolog.InfoLevel for info only messages etc */

	Global = log.Logger
}

// NewLogger creates a new logger with the given options.
func NewLogger(opts ...*Options) zerolog.Logger {
	var opt = defaultOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	// We don't want callers in the log output if it isn't enabled.
	var config = logConfig

	if !opt.Callers {
		config.FormatCaller = func(i interface{}) string {
			return ""
		}
	}

	/* make logger :3 heh */

	logger := log.Output(config).With().Logger()
	if opt.MinimumLevel != zerolog.NoLevel {
		logger = logger.Level(opt.MinimumLevel)
	}

	return logger
}
