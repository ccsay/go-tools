package jsonWebToken

import (
	"github.com/liuchonglin/go-utils"
)

var tokenConfig *TokenConfig
// token配置
type TokenConfig struct {
	// 发行者
	Issuer string `json:"issuer" yaml:"issuer"`
	// 过期时间
	ExpireTime int64 `json:"expireTime" yaml:"expireTime"`
	// 刷新过期时间
	RefreshExpireTime int64 `json:"refreshExpireTime" yaml:"refreshExpireTime"`
	// 忽略方法列表
	IgnoreMethods []string `json:"ignoreMethods" yaml:"ignoreMethods"`
}

func GetTokenConfig() *TokenConfig {
	if tokenConfig == nil {
		tokenConfig = &TokenConfig{}
		tokenConfig.defaultValue()
	}
	return tokenConfig
}

func (t *TokenConfig) defaultValue() {
	if utils.IsEmpty(t.Issuer) {
		t.Issuer = "go-tools"
	}
	if t.ExpireTime == 0 {
		t.ExpireTime = 3600
	}
	if t.RefreshExpireTime == 0 {
		t.RefreshExpireTime = 3600
	}
}
