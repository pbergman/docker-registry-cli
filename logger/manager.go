package logger

import (
	"github.com/pbergman/logger"
	"github.com/pbergman/logger/handlers"
	"github.com/pbergman/logger/processors"
)

var Logger *logger.Logger

func init() {
	Logger = logger.NewLogger("drc")
}

func SetHandler(verbose bool) {
	if verbose {
		processor := processors.NewTraceProcessor(logger.WARNING)
		Logger.AddProcessor(processor.Process)
		Logger.AddHandler(handlers.NewStdoutHandler(logger.DEBUG))
	} else {
		Logger.AddHandler(handlers.NewStdoutHandler(logger.INFO))
	}
}
