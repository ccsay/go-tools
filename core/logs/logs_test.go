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

// 日志工具类
package logs

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/liuchonglin/go-tools/core/common"
)

var config = &LogsConfig{
	PrintLine:     true,
	OutputConsole: true,
	OutputFile:    true,
}

func TestInit(t *testing.T) {
	tests := []struct {
		name         string
		wantErr      bool
		updateConfig func()
	}{
		{
			name:    "all",
			wantErr: false,
			updateConfig: func() {

			},
		}, {
			name:    "logsConfig nil",
			wantErr: false,
			updateConfig: func() {
				logsConfig = nil
			},
		},
	}
	for _, tt := range tests {
		tt.updateConfig()
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(config); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		ctx    context.Context
		logTag string
	}
	tests := []struct {
		name string
		args args
		want *Logs
	}{
		{
			name: "ok",
			args: args{ctx: context.Background(), logTag: logTag},
			want: &Logs{ctx: context.Background(), logTag: logTag},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.ctx, tt.args.logTag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogs_Debug(t *testing.T) {
	type fields struct {
		ctx    context.Context
		logTag string
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "all",
			fields: fields{ctx: context.Background(), logTag: logTag},
			args:   args{format: "传入参数a=%v b=%v", args: []interface{}{1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logs{
				ctx:    tt.fields.ctx,
				logTag: tt.fields.logTag,
			}
			l.Debug(tt.args.format, tt.args.args...)
		})
	}
}

func TestLogs_Info(t *testing.T) {
	type fields struct {
		ctx    context.Context
		logTag string
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "all",
			fields: fields{ctx: context.Background(), logTag: logTag},
			args:   args{format: "传入参数a=%v b=%v", args: []interface{}{1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logs{
				ctx:    tt.fields.ctx,
				logTag: tt.fields.logTag,
			}
			l.Info(tt.args.format, tt.args.args...)
		})
	}
}

func TestLogs_Warn(t *testing.T) {
	type fields struct {
		ctx    context.Context
		logTag string
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "all",
			fields: fields{ctx: context.Background(), logTag: logTag},
			args:   args{format: "传入参数a=%v b=%v", args: []interface{}{1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logs{
				ctx:    tt.fields.ctx,
				logTag: tt.fields.logTag,
			}
			l.Warn(tt.args.format, tt.args.args...)
		})
	}
}

func TestLogs_Error(t *testing.T) {
	type fields struct {
		ctx    context.Context
		logTag string
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "all",
			fields: fields{ctx: context.Background(), logTag: logTag},
			args:   args{format: "传入参数a=%v b=%v", args: []interface{}{1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logs{
				ctx:    tt.fields.ctx,
				logTag: tt.fields.logTag,
			}
			l.Error(tt.args.format, tt.args.args...)
		})
	}
}

func Test_ctxValue(t *testing.T) {
	ctxString := context.WithValue(context.Background(), common.TraceIdKey, "123456")
	ctxInt := context.WithValue(context.Background(), common.TraceIdKey, 123456)
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{ctx: ctxString, key: common.TraceIdKey},
			want: "123456",
		}, {
			name: "ctx is nil",
			args: args{ctx: nil, key: common.TraceIdKey},
			want: "",
		}, {
			name: "value int",
			args: args{ctx: ctxInt, key: common.TraceIdKey},
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ctxValue(tt.args.ctx, tt.args.key); got != tt.want {
				t.Errorf("ctxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertLogLevel(t *testing.T) {
	type args struct {
		levelStr string
	}
	tests := []struct {
		name string
		args args
		want zapcore.Level
	}{
		{
			name: "default",
			args: args{
				levelStr: "",
			},
			want: zap.InfoLevel,
		}, {
			name: "debug",
			args: args{
				levelStr: DEBUG,
			},
			want: zap.DebugLevel,
		}, {
			name: "info",
			args: args{
				levelStr: INFO,
			},
			want: zap.InfoLevel,
		}, {
			name: "warn",
			args: args{
				levelStr: WARN,
			},
			want: zap.WarnLevel,
		}, {
			name: "error",
			args: args{
				levelStr: ERROR,
			},
			want: zap.ErrorLevel,
		}, {
			name: "debug",
			args: args{
				levelStr: FATAL,
			},
			want: zap.FatalLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertLogLevel(tt.args.levelStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	if err := Init(config); err != nil {
		panic(fmt.Sprintf("Init() error = %v", err))
	}
	os.Exit(m.Run())
}
