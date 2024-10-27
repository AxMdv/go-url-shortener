package main

import (
	"os"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
)

func Test_formatValue(t *testing.T) {
	type args struct {
		buildData string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive #1",
			args: args{buildData: "v1.0.0"},
			want: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatValue(tt.args.buildData); got != tt.want {
				t.Errorf("formatValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	_ = config.ParseOptions()
	os.Exit(m.Run())
}
