package logger

import "go.uber.org/zap"

var Log *zap.Logger = zap.NewNop()

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
