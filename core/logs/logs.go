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

// 日志工具类
package logs

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"strings"
	"os"
	"github.com/liuchonglin/go-utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/liuchonglin/go-tools/core/common"
)

const (
	logTagKey      = "logTag"
	serviceNameKey = "serviceName"
	logTag         = "core.logs"
)

type Logs struct {
	ctx    context.Context
	logTag string
	logger *zap.Logger
}

type Logger struct {
	Logger *zap.Logger
}

func NewLogger(logsConfig *LogsConfig) (*Logger, error) {
	if logsConfig == nil {
		logsConfig = &LogsConfig{FileLogger: &lumberjack.Logger{}}
	}
	// 配置默认值
	if err := logsConfig.defaultValue(); err != nil {
		return nil, err
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(convertLogLevel(logsConfig.Level))

	// 打印到控制台和文件
	var ws zapcore.WriteSyncer
	writeConsole := zapcore.AddSync(os.Stdout)
	writeFile := zapcore.AddSync(logsConfig.FileLogger)
	if logsConfig.OutputConsole && logsConfig.OutputFile {
		ws = zapcore.NewMultiWriteSyncer(writeConsole, writeFile)
	} else if logsConfig.OutputFile {
		ws = writeFile
	} else {
		ws = writeConsole
	}

	// 创建一个Core，将日志写入WriteSyncer
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), ws, atomicLevel)
	// 设置初始化字段
	filed := zap.Fields(zap.String(serviceNameKey, logsConfig.ServiceName))
	var logger *zap.Logger
	// 开启文件及行号
	if logsConfig.PrintLine {
		caller := zap.AddCaller()
		// 跳过的调用者数量，如果不设置该值会打印本文件的行号
		callerSkip := zap.AddCallerSkip(1)
		logger = zap.New(core, filed, caller, callerSkip)
	} else {
		logger = zap.New(core, filed)
	}
	return &Logger{Logger: logger}, nil
}

func New(ctx context.Context, logTag string, zapLogger *zap.Logger) *Logs {
	var logger *Logger
	if utils.IsEmpty(zapLogger) {
		var err error
		logger, err = NewLogger(nil)
		if err != nil {
			panic(err)
		}
		zapLogger = logger.Logger
	}
	logs := &Logs{ctx: ctx, logTag: logTag, logger: zapLogger}
	return logs
}

// 最低等级的，主要用于开发过程中打印一些运行/调试信息，不允许生产环境打开debug级别
func (l *Logs) Debug(format string, args ...interface{}) {
	msg := format
	if len(args) != 0 {
		// 注意：如果args为空，fmt.Sprintf()会格式化错误，多出 %!(EXTRA []interface {}=[]) 字符串
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Debug(msg, zap.String(common.TraceIdKey, ctxValue(l.ctx, common.TraceIdKey)),
		zap.String(logTagKey, l.logTag))
}

// 打印一些你感兴趣的或者重要的信息，这个可以用于生产环境中输出程序运行的一些重要信息
func (l *Logs) Info(format string, args ...interface{}) {
	msg := format
	if len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Info(msg, zap.String(common.TraceIdKey, ctxValue(l.ctx, common.TraceIdKey)),
		zap.String(logTagKey, l.logTag))
}

// 表明会出现潜在错误的情形，有些信息不是错误信息，但是也要给程序员的一些提示
func (l *Logs) Warn(format string, args ...interface{}) {
	msg := format
	if len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Warn(msg, zap.String(common.TraceIdKey, ctxValue(l.ctx, common.TraceIdKey)),
		zap.String(logTagKey, l.logTag))
}

// 指出虽然发生错误事件，但仍然不影响系统的继续运行。打印错误和异常信息
func (l *Logs) Error(format string, args ...interface{}) {
	msg := format
	if len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Error(msg, zap.String(common.TraceIdKey, ctxValue(l.ctx, common.TraceIdKey)),
		zap.String(logTagKey, l.logTag))
}

// 指出每个严重的错误事件将会导致应用程序的退出。这个级别比较高了。重大错误
func (l *Logs) Fatal(format string, args ...interface{}) {
	msg := format
	if len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Fatal(msg, zap.String(common.TraceIdKey, ctxValue(l.ctx, common.TraceIdKey)),
		zap.String(logTagKey, l.logTag))
}

// 根据 key 从上下文(context)获取值
func ctxValue(ctx context.Context, key string) string {
	var value string
	if ctx == nil {
		return value
	}

	if v, ok := ctx.Value(key).(string); ok {
		value = v
	}

	if v, ok := ctx.Value(key).(int); ok {
		value = strconv.Itoa(v)
	}
	return value
}

// 把字符串转换为日志级别（数字）
func convertLogLevel(levelStr string) zapcore.Level {
	// 不区分大小写
	levelStr = strings.ToLower(levelStr)
	var level zapcore.Level
	switch levelStr {
	case DEBUG:
		level = zap.DebugLevel
	case INFO:
		level = zap.InfoLevel
	case WARN:
		level = zap.WarnLevel
	case ERROR:
		level = zap.ErrorLevel
	case FATAL:
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	return level
}
