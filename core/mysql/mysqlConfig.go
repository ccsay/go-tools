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

// mysql配置
package mysql

import (
	"fmt"
	"github.com/liuchonglin/go-utils"
)

var mysqlConfig *MysqlConfig

// mysql 配置
type MysqlConfig struct {
	// 用户名
	Username string `json:"username" yaml:"username"`
	// 密码
	Password string `json:"password" yaml:"password"`
	// IP地址
	Address string `json:"address" yaml:"address"`
	// 数据库
	Database string `json:"database" yaml:"database"`
	// 时区
	Timezone string `json:"timezone" yaml:"timezone"`
	// 最大连接数
	MaxConn int `json:"maxConn" yaml:"maxConn"`
	// 最大超时事
	MaxIdle int `json:"maxIdle" yaml:"maxIdle"`
	// 连接最大存活时间
	MaxLifeTime int64 `json:"maxLifeTime" yaml:"maxLifeTime"`
	// 数据库源
	Source string `json:"source"`
	// 是否显示SQL
	ShowSql bool `json:"showSql" yaml:"showSql"`
}

func GetMysqlConfig() (*MysqlConfig, error) {
	if mysqlConfig == nil {
		mysqlConfig = &MysqlConfig{}
		if err := mysqlConfig.defaultValue(); err != nil {
			return nil, err
		}
	}
	return mysqlConfig, nil
}

func (m *MysqlConfig) defaultValue() error {
	if utils.IsEmpty(m.Database) {
		return fmt.Errorf("mysql 数据库 [ datebase ] 不能为空")
	}
	if utils.IsEmpty(m.Username) {
		m.Username = "root"
	}
	if utils.IsEmpty(m.Password) {
		m.Password = "root"
	}
	if utils.IsEmpty(m.Address) {
		m.Address = "localhost:3306"
	}
	if utils.IsEmpty(m.Timezone) {
		m.Timezone = "Local"
	}
	if m.MaxIdle == 0 {
		m.MaxIdle = 16
	}
	if m.MaxConn == 0 {
		m.MaxConn = 16
	}
	if m.MaxLifeTime == 0 {
		m.MaxLifeTime = 600
	}

	m.Source = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=%s",
		m.Username,
		m.Password,
		m.Address,
		m.Database,
		m.Timezone)
	return nil
}
