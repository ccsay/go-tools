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

// 日志配置
package logs

import (
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestGetLogsConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLogsConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogsConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLogsConfig_defaultValue(t *testing.T) {
	type fields struct {
		FileLogger    *lumberjack.Logger
		PrintLine     bool
		Level         string
		ServiceName   string
		OutputConsole bool
		OutputFile    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				FileLogger:    nil,
				PrintLine:     false,
				Level:         "",
				ServiceName:   "",
				OutputConsole: true,
				OutputFile:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogsConfig{
				FileLogger:    tt.fields.FileLogger,
				PrintLine:     tt.fields.PrintLine,
				Level:         tt.fields.Level,
				ServiceName:   tt.fields.ServiceName,
				OutputConsole: tt.fields.OutputConsole,
				OutputFile:    tt.fields.OutputFile,
			}
			if err := l.defaultValue(); (err != nil) != tt.wantErr {
				t.Errorf("LogsConfig.defaultValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
