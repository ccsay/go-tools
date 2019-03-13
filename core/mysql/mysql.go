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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type MysqlDB struct {
	DB *gorm.DB
}

func NewMysqlDB(mysqlConfig *MysqlConfig) (*MysqlDB, error) {
	if mysqlConfig == nil {
		mysqlConfig = &MysqlConfig{}
	}
	if err := mysqlConfig.defaultValue(); err != nil {
		return nil, err
	}
	var err error
	db, err := gorm.Open("mysql", mysqlConfig.Source)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(mysqlConfig.MaxIdle)
	db.DB().SetMaxOpenConns(mysqlConfig.MaxConn)
	db.DB().SetConnMaxLifetime(time.Duration(mysqlConfig.MaxLifeTime) * time.Second)
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	// 设置是否打印SQL
	db.LogMode(mysqlConfig.ShowSql)
	return &MysqlDB{DB: db}, nil
}
