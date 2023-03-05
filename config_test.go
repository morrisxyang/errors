package errors

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCfg(t *testing.T) {
	tests := []struct {
		cfg  *Config
		want *Config
	}{
		{
			nil,
			cfg,
		},
		{
			&Config{
				StackDepth:          100,
				ErrorConnectionFlag: ":",
			},
			&Config{
				StackDepth:          100,
				ErrorConnectionFlag: ":",
			},
		},
		{
			&Config{
				StackDepth: 100,
			},
			&Config{
				StackDepth: 100,
			},
		},
	}
	for i, tt := range tests {
		SetCfg(tt.cfg)
		got := cfg
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func TestGetCfg(t *testing.T) {
	expected := defaultCfg
	ResetCfg()
	// call GetCfg
	c := GetCfg()

	// check if returned value is the expected one
	assert.Equal(t, expected, c)
}

func TestResetCfg(t *testing.T) {
	SetCfg(&Config{
		StackDepth: 100,
	})
	assert.NotEqual(t, defaultCfg, cfg)

	ResetCfg()
	// check if config was reset to default values
	assert.Equal(t, defaultCfg, cfg)
}
