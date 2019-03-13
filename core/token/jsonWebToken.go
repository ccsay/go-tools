package jsonWebToken

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	logTag       = "core.token"
	TokenKey     = "token"
	TokenDataKey = "tokenData"
)

// JWT签名结构
type jsonWebToken struct {
	SigningKey []byte
}

// 自定义载荷 必须继承 jwt.StandardClaims
type customClaims struct {
	Data map[string]interface{}
	jwt.StandardClaims
}

func New(tokenConfig *TokenConfig) *jsonWebToken {
	if tokenConfig == nil {
		tokenConfig = &TokenConfig{}
		tokenConfig.defaultValue()
	}
	return &jsonWebToken{}
}

// 创建token
func (j *jsonWebToken) CreateToken(data map[string]interface{}) (string, error) {
	claims := &customClaims{Data: data, StandardClaims: jwt.StandardClaims{
		//签名生效时间
		NotBefore: int64(time.Now().Unix() - 1000),
		//签名过期时间 1小时
		ExpiresAt: int64(time.Now().Unix() + GetTokenConfig().ExpireTime),
		//签名发行者
		Issuer: GetTokenConfig().Issuer,
	}}
	return j.createToken(claims)
}

//解析token
func (j *jsonWebToken) ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
			//logs.New(context.Background(), logTag).Info("解析token成功: %v,token = %v", claims, tokenString)
			return claims.Data, nil
		}
	}
	//logs.New(context.Background(), logTag).Warn("解析token错误: %v, token = %v", err, tokenString)
	return nil, err
}

//更新Token
func (j *jsonWebToken) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
			//tempExpires := claims.ExpiresAt
			//logs.New(context.Background(), logTag).Info("token更新前 expires = %v,token = %v", tempExpires, tokenString)
			jwt.TimeFunc = time.Now
			refreshExpireTime := time.Duration(GetTokenConfig().RefreshExpireTime)
			claims.ExpiresAt = time.Now().Add(refreshExpireTime * time.Second).Unix()
			//logs.New(context.Background(), logTag).Info("token更新后 expires = %v,token = %v", claims.ExpiresAt, tokenString)
			return j.createToken(claims)
		}
	}
	//logs.New(context.Background(), logTag).Warn("解析token错误: %v, token = %v", err, tokenString)
	return "", err
}

// 创建token
func (j *jsonWebToken) createToken(claims *customClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.SigningKey)
	if err != nil {
		//logs.New(context.Background(), logTag).Error("创建token错误: %v, claims = %v", err, claims)
		return "", err
	}
	//logs.New(context.Background(), logTag).Info("创建token成功 token = %v", tokenString)
	return tokenString, nil
}
