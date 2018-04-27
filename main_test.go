package main

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestReadKeys(t *testing.T) {
	type args struct {
		configfile string
	}
	tests := []struct {
		name     string
		args     args
		wantKeys map[string]string
		fixRand  bool
	}{
		{
			name: "simple sample file",
			args: args{configfile: "test/keys.toml"},
			wantKeys: map[string]string{
				"asdf":  "sample",
				"asdf2": "sample2",
			},
		},
		{
			name: "random generated",
			args: args{configfile: ""},
			wantKeys: map[string]string{
				"cUbYhiZzKa": "GENERATED",
			},
			fixRand: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fixRand {
				rand.Seed(0)
			}
			if gotKeys := ReadKeys(tt.args.configfile); !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("ReadKeys() = %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func Test_checkKey(t *testing.T) {
	type args struct {
		key string
	}
	keys["asdf"] = "ASDF"
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "existing key",
			args: args{"asdf"},
			want: true,
		},
		{
			name: "not existing key",
			args: args{"asdf2"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkKey(tt.args.key); got != tt.want {
				t.Errorf("checkKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
