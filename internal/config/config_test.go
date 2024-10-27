package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	opt := ParseOptions()
	assert.Equal(t, ":8080", opt.RunAddr)
	assert.Equal(t, "http://localhost:8080", opt.ResponseResultAddr)
	assert.Equal(t, "/tmp/short-url-db.json", opt.FileStorage)
	assert.Equal(t, "", opt.DataBaseDSN)
}
