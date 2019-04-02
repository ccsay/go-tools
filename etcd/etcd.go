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

// etcd工具类
package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
	"github.com/liuchonglin/go-utils"
)

type etcd struct {
	EtcdClient *clientv3.Client
}

var (
	etcdClientIsNilError = errors.New("etcd client is nil")
	keyEmptyError        = errors.New("etcd key is empty")
	valueNotJson         = errors.New("'value' is not a json")
)

func NewEtcd(etcdConfig *EtcdConfig) (e *etcd, err error) {
	if etcdConfig != nil {
		etcdConfig = &EtcdConfig{}
	}
	etcdConfig.defaultValue()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdConfig.Address},
		DialTimeout: time.Duration(etcdConfig.Timeout) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd 初始化失败: %v", err)
	}
	//etcd超时控制, 设置ContextTimeout超时
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(GetEtcdConfig().ContextTimeout)*time.Second)
	_, err = etcdClient.Get(ctx, "init_get_test_key")
	if err != nil {
		return nil, fmt.Errorf("etcd 初始化失败: %v", err)
	}
	//操作完毕，取消超时控制
	cancel()
	return &etcd{EtcdClient: etcdClient}, nil
}

func (e *etcd) Get(key string, config interface{}) error {
	if e.EtcdClient == nil {
		return etcdClientIsNilError
	}
	// 定义key  /公司名/项目名/before（接口服务） or after（后台服务）/配置名称
	if key = formatKey(key); utils.IsEmpty(key) {
		return keyEmptyError
	}
	if err := utils.CheckPointer(config); err != nil {
		return err
	}
	value, err := get(e.EtcdClient, key)
	if err != nil {
		return nil
	}
	return json.Unmarshal(value, config)
}

// 通过key 从etcd中获取value
func get(etcdClient *clientv3.Client, key string) (value []byte, err error) {
	//etcd超时控制, 设置ContextTimeout超时
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(GetEtcdConfig().ContextTimeout)*time.Second)
	resp, err := etcdClient.Get(ctx, key)
	//操作完毕，取消超时控制
	cancel()
	if err != nil || len(resp.Kvs) == 0 {
		return nil, err
	}
	return resp.Kvs[0].Value, nil
}

func (e *etcd) Watch(key string, f func(event *clientv3.Event)) error {
	if e.EtcdClient == nil {
		return etcdClientIsNilError
	}
	if key = formatKey(key); utils.IsEmpty(key) {
		return keyEmptyError
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("got panic in watch: %+v", r)
			}
			watch(e.EtcdClient, key, f)
		}()
	}()
	return nil
}

func watch(etcdClient *clientv3.Client, key string, f func(event *clientv3.Event)) {
	for {
		watchChan := etcdClient.Watch(context.Background(), key)
		for wResp := range watchChan {
			for _, ev := range wResp.Events {
				f(ev)
			}
		}
	}
}

func (e *etcd) Put(key string, value string) error {
	if e.EtcdClient == nil {
		return etcdClientIsNilError
	}
	if key = formatKey(key); utils.IsEmpty(key) {
		return keyEmptyError
	}
	if flag, err := validateJson(value); err != nil || !flag {
		return valueNotJson
	}
	return put(e.EtcdClient, key, value)
}

func put(etcdClient *clientv3.Client, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(GetEtcdConfig().ContextTimeout)*time.Second)
	_, err := etcdClient.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func (e *etcd) Delete(key string) error {
	if e.EtcdClient == nil {
		return etcdClientIsNilError
	}
	if key = formatKey(key); utils.IsEmpty(key) {
		return keyEmptyError
	}
	return delete(e.EtcdClient, key)
}

func delete(etcdClient *clientv3.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(GetEtcdConfig().ContextTimeout)*time.Second)
	_, err := etcdClient.Delete(ctx, key)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

// 格式化key
func formatKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	key = strings.TrimSuffix(key, "/")
	if strings.Contains(key, "//") {
		return ""
	}
	return key
}

// 验证是否是json
func validateJson(data string) (bool, error) {
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		return false, err
	}
	return true, nil
}
