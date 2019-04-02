package filter

import (
	"github.com/liuchonglin/go-tools/token"
	"net/http"
	"fmt"
	"strings"
	"github.com/liuchonglin/go-utils"
	"context"
	"github.com/liuchonglin/go-tools/result"
	"time"
	"github.com/liuchonglin/go-tools/common"
	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

const (
	logTag = "core.filter"
)

var limiter = rate.NewLimiter(40000, 20000)

// 限流过滤器
func RateFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// AllowN标识在时间now的时候，n个事件是否可以同时发生(也意思就是now的时候是否可以从令牌池中取n个令牌)。
		// 如果你需要在事件超出频率的时候丢弃或跳过事件
		if !limiter.AllowN(time.Now(), 1) {
			c.JSON(http.StatusOK, result.NewError(common.SystemBusy, common.SystemBusyMessage))
			// 终止
			c.Abort()
			return
		}
		// 执行下一个中间件
		c.Next()
	}
}

// 上下文处理
func ContextHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(common.TraceIdKey)
		if utils.IsEmpty(traceId) || len(traceId) != 36 {
			//logs.New(context.Background(), logTag).Error(common.ParamErrorMessage, "traceId")
			c.JSON(http.StatusOK, result.NewError(common.ParamError, fmt.Sprintf(common.ParamErrorMessage, "traceId")))
			c.Abort()
			return
		}
		ctx := context.WithValue(context.Background(), common.TraceIdKey, traceId)
		c.Set(common.CtxKey, ctx)
		c.Next()
	}
}

// JWT处理
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if isIgnorePath(c) {
			c.Next()
			return
		}

		token := c.GetHeader(jsonWebToken.TokenKey)
		if utils.IsEmpty(token) {
			//logs.New(context.Background(), logTag).Warn(common.ParamErrorMessage, "token is empty")
			c.JSON(http.StatusOK, result.NewError(common.ParamError, fmt.Sprintf(common.ParamErrorMessage, "token is empty")))
			c.Abort()
			return
		}

		// 验证Token
		j := jsonWebToken.New(jsonWebToken.GetTokenConfig())
		data, err := j.ParseToken(token)
		if err != nil {
			//logs.New(context.Background(), logTag).Error(common.ParamErrorMessage, err)
			c.JSON(http.StatusOK, result.NewError(common.ParamError, fmt.Sprintf(common.ParamErrorMessage, err)))
			c.Abort()
			return
		}
		// 把Token中的数据放入上下文
		ctx := context.WithValue(context.Background(), jsonWebToken.TokenDataKey, data)
		c.Set(common.CtxKey, ctx)
		c.Next()
	}
}

// 是否是忽略路径
func isIgnorePath(c *gin.Context) bool {
	methods := jsonWebToken.GetTokenConfig().IgnoreMethods
	if methods != nil && len(methods) > 0 {
		path := c.Request.URL.Path
		for _, method := range methods {
			if strings.Index(path, method) != -1 {
				return true
			}
		}
	}
	return false
}
