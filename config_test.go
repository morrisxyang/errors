package errors

import (
	"reflect"
	"testing"
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
				Depth:               100,
				ErrorConnectionFlag: ":",
			},
			&Config{
				Depth:               100,
				ErrorConnectionFlag: ":",
			},
		},
		{
			&Config{
				Depth: 100,
			},
			&Config{
				Depth: 100,
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
