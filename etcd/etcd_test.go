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
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/liuchonglin/go-utils/timeutil"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

const (
	testConfigKey = "/liuchonglin/test/before/product"
)

type ProductConfig struct {
	ProductId int64     `json:"productId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Status    int8      `json:"status"`
	Total     int       `json:"total"`
}

var etcdClient *etcd

var config = &EtcdConfig{
	Address:        "localhost:2379",
	Timeout:        3,
	ContextTimeout: 10,
}

func TestNewEtcd(t *testing.T) {
	type args struct {
		etcdConfig *EtcdConfig
	}
	tests := []struct {
		name     string
		args     args
		wantEtcd *etcd
		wantErr  bool
	}{
		{
			name: "all",
			args: args{
				etcdConfig: &EtcdConfig{
					Address:        "localhost:2379",
					Timeout:        3,
					ContextTimeout: 10,
				},
			},
			wantErr: false,
		}, {
			name: "address err",
			args: args{
				etcdConfig: &EtcdConfig{
					Address:        "",
					Timeout:        3,
					ContextTimeout: 10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEtcd(tt.args.etcdConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEtcd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestPut(t *testing.T) {
	startTime, _ := timeutil.StringToTime("2019-01-10 12:00:00", timeutil.FormatTime)
	endTime, _ := timeutil.StringToTime("2019-01-11 12:00:00", timeutil.FormatTime)
	configMap := make(map[int]ProductConfig)
	configMap[1] = ProductConfig{ProductId: 1,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    0,
		Total:     100,
	}
	configMap[2] = ProductConfig{ProductId: 2,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    0,
		Total:     1000,
	}

	configMap[3] = ProductConfig{ProductId: 3,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    0,
		Total:     2000,
	}

	jsonData, err := json.Marshal(configMap)
	if err != nil {
		t.Errorf("json 序列化错误: %v", err)
		return
	}

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				key:   testConfigKey,
				value: string(jsonData),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := etcdClient.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	configMap := make(map[int]*ProductConfig)

	type args struct {
		key    string
		config interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "all",
			args:    args{key: testConfigKey, config: &configMap},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := etcdClient.Get(tt.args.key, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args.config)
		})
	}
}

func TestWatch(t *testing.T) {
	productConfig := &ProductConfig{}
	type args struct {
		key string
		f   func(event *clientv3.Event)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				key: testConfigKey,
				f: func(event *clientv3.Event) {
					if event.Type == mvccpb.DELETE {
						fmt.Printf("key[%s] 's config deleted", event.Kv.Key)
					}
					if event.Type == mvccpb.PUT && string(event.Kv.Key) == testConfigKey {
						err := json.Unmarshal(event.Kv.Value, productConfig)
						if err != nil {
							fmt.Printf("key [%s], Unmarshal[%s], err:%v", event.Kv.Key, event.Kv.Value, err)
						}
					}
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := etcdClient.Watch(tt.args.key, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	time.Sleep(1 * time.Second)
}

func TestDelete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{key: testConfigKey},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := etcdClient.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_formatKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				key: "/caitong",
			},
			want: "/caitong",
		}, {
			name: "prefix",
			args: args{
				key: "caitong",
			},
			want: "/caitong",
		}, {
			name: "Suffix",
			args: args{
				key: "/caitong/",
			},
			want: "/caitong",
		}, {
			name: "double",
			args: args{
				key: "/caitong//spikeProxy/",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatKey(tt.args.key); got != tt.want {
				t.Errorf("formatKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateJson(t *testing.T) {
	startTime, _ := timeutil.StringToTime("2019-01-10 12:00:00", timeutil.FormatTime)
	endTime, _ := timeutil.StringToTime("2019-01-11 12:00:00", timeutil.FormatTime)
	productConfig := ProductConfig{ProductId: 10,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    0,
		Total:     100,
	}
	jsonData, err := json.Marshal(productConfig)
	if err != nil {
		t.Errorf("json 序列化错误：%v", err)
		return
	}

	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				data: string(jsonData),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				data: "123",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateJson(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	var err error
	etcdClient, err = NewEtcd(config)
	if err != nil {
		panic(fmt.Sprintf("Init() error = %v", err))
	}

	os.Exit(m.Run())
}
