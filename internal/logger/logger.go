package logger

import (
	"fmt"
	"log"
	"strings"
)

// New returns a new custom logger
func New(prefix string) (*log.Logger, error) {
	prefix = strings.ToUpper(prefix)

	switch prefix {
	case "INFO", "DEBUG", "ERROR":
	default:
		return nil, fmt.Errorf("invalid logger prefix: %s", prefix)
	}

	logger := log.Default()
	logger.SetPrefix(prefix + ": ")

	return logger, nil
}
