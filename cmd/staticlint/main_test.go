package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = getStandarStaticAnalyzers()
	_ = getStaticcheckAnalyzers()
	_ = getQuickfixAnalyzers()
	_ = getCommunityAnalyzers()

	os.Exit(m.Run())
}
