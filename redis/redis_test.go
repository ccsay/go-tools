// Copyright 2019 go-tools Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// redis工具类
package redis

import (
	"testing"
	"os"
)

var config = &RedisConfig{
	Address:     "localhost:6379",
	MaxIdle:     16,
	MaxActive:   100,
	IdleTimeout: 300,
	Password:    "liu5522112",
}

var redisTool *Redis

func TestNewRedis(t *testing.T) {
	type args struct {
		redisConfig *RedisConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all",
			args: args{
				redisConfig: &RedisConfig{
					Address:     "localhost:6379",
					MaxIdle:     16,
					MaxActive:   100,
					IdleTimeout: 300,
					Password:    "liu5522112",
				},
			},
			wantErr: false,
		}, {
			name: "redisConfig nil",
			args: args{
				redisConfig: nil,
			},
			wantErr: false,
		}, {
			name: "password nil",
			args: args{
				redisConfig: &RedisConfig{
					Address:     "localhost:6379",
					MaxIdle:     16,
					MaxActive:   100,
					IdleTimeout: 300,
					Password:    "",
				},
			},
			wantErr: false,
		}, {
			name:    "link err",
			wantErr: true,
			args: args{
				redisConfig: &RedisConfig{
					Address:     "",
					MaxIdle:     16,
					MaxActive:   100,
					IdleTimeout: 300,
					Password:    "liu5522112",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRedis(tt.args.redisConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name:    "all",
			args:    args{key: "name", value: "xiaoliu"},
			wantErr: false,
		}, {
			name:    "key nil",
			args:    args{key: "", value: "xiaoliu"},
			wantErr: true,
		}, {
			name:    "value nil",
			args:    args{key: "name", value: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.mock != nil {
			tt.mock()
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := redisTool.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
		wantErr   bool
	}{
		{
			name:      "all",
			args:      args{key: "name"},
			wantValue: "xiaoliu",
			wantErr:   false,
		}, {
			name:      "key nil",
			args:      args{key: ""},
			wantValue: "",
			wantErr:   true,
		}, {
			name:      "key not exist",
			args:      args{key: "name12345"},
			wantValue: "",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := redisTool.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValue != tt.wantValue {
				t.Errorf("Get() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestDel(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "all",
			args:    args{key: "name"},
			wantErr: false,
		}, {
			name:    "key nil",
			args:    args{key: ""},
			wantErr: true,
		}, {
			name:    "key not exist",
			args:    args{key: "name123456"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := redisTool.Del(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetExpire(t *testing.T) {
	type args struct {
		key   string
		value string
		ex    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "all",
			args:    args{key: "age", value: "20", ex: 60},
			wantErr: false,
		}, {
			name:    "key nil",
			args:    args{key: "", value: "20", ex: 60},
			wantErr: true,
		}, {
			name:    "value nil",
			args:    args{key: "age", value: "", ex: 60},
			wantErr: true,
		}, {
			name:    "ex is zero",
			args:    args{key: "age", value: "20", ex: 0},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := redisTool.SetExpire(tt.args.key, tt.args.value, tt.args.ex); (err != nil) != tt.wantErr {
				t.Errorf("SetExpire() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMain(m *testing.M) {
	var err error
	redisTool, err = NewRedis(config)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
