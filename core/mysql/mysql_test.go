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

// mysql工具类
package mysql

import (
	"testing"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var config = &MysqlConfig{
	Username:    "root",
	Password:    "5522112",
	Address:     "localhost:3306",
	Database:    "base_cms",
	Timezone:    "Local",
	MaxConn:     10,
	MaxIdle:     10,
	MaxLifeTime: 100,
	ShowSql:     true,
}

func TestInit(t *testing.T) {
	type args struct {
		mysqlConfig *MysqlConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				mysqlConfig: config,
			},
			wantErr: false,
		}, {
			name: "nil",
			args: args{
				mysqlConfig: nil,
			},
			wantErr: true,
		}, {
			name: "empty",
			args: args{
				mysqlConfig: &MysqlConfig{Database: "base_cms"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.args.mysqlConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMysql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
