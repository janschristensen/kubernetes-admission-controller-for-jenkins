package main

import (
	"reflect"
	"testing"
)

func Test_readConf(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     []string
		wantErr  bool
	}{
		{"Default", "image_conf.yaml", []string{"Default"}, false},
		{"Allow.Empty", "image_conf.yaml", []string{"Allow", "Empty"}, false},
		{"Allow.Registries", "image_conf.yaml", []string{"Allow", "Registries"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readConf(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("readConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			image := reflect.Indirect(reflect.ValueOf(&got).Elem())
			var currentValue reflect.Value
			currentValue = image
			for _, ttt := range tt.want {
				currentValue = currentValue.FieldByName(ttt)
			}
			if currentValue.Interface() == "" {
				t.Errorf("readConf() = %v, want non-empty", currentValue.Interface())
			}
		})
	}
}
