// Package logger initializes logger.
package logger

import "go.uber.org/zap"

// Log is the singleton logger.
var Log *zap.Logger = zap.NewNop()

// InitLogger() inits Log according to logging level.
func InitLogger() error {

	lvl, err := zap.ParseAtomicLevel("info")
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
