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

// etcd配置
package etcd

import (
	"github.com/liuchonglin/go-utils"
)

var etcdConfig *EtcdConfig

// etcd 配置
type EtcdConfig struct {
	// 连接地址
	Address string `json:"address" yaml:"address"`
	// 连接超时时间
	Timeout int64 `json:"timeout" yaml:"timeout"`
	// 设置ContextTimeout超时
	ContextTimeout int64 `json:"contextTimeout" yaml:"contextTimeout"`
}

func GetEtcdConfig() *EtcdConfig {
	if etcdConfig == nil {
		etcdConfig = &EtcdConfig{}
		etcdConfig.defaultValue()
	}
	return etcdConfig
}

func (e *EtcdConfig) defaultValue() {
	if utils.IsEmpty(e.Address) {
		e.Address = "localhost:2379"
	}
	if e.Timeout == 0 {
		e.Timeout = 5
	}
	if e.ContextTimeout == 0 {
		e.ContextTimeout = 10
	}
}
