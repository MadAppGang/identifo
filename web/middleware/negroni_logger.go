package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/negroni"
)

func NewNegroniLogger(name string) *negroni.Logger {
	logger := negroni.NewLogger()
	logger.ALogger = log.New(os.Stdout, fmt.Sprintf("[ %s ]: ", name), 0)
	logger.SetFormat(negroni.LoggerDefaultFormat)
	return logger
}
