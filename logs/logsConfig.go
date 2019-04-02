package logs

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"runtime"
	"path"
	"github.com/liuchonglin/go-utils"
	"path/filepath"
	"os"
	"github.com/liuchonglin/go-utils/fileutil"
)

const (
	logsFilePath   = "logs"
	logsFileName   = "access.log"
	logsFileSuffix = ".log"
	DEBUG          = "debug"
	INFO           = "info"
	WARN           = "warn"
	ERROR          = "error"
	FATAL          = "fatal"
)

var logsConfig *LogsConfig

// logs 配置
type LogsConfig struct {
	// 文件滚动配置
	FileLogger *lumberjack.Logger `json:"fileLogger" yaml:"fileLogger"`
	// 是否打印行数
	// 默认值：false
	PrintLine bool `json:"printLine" yaml:"printLine"`
	// 日志级别
	// 可选值：debug|info|warn|error|fatal
	Level string `json:"level" yaml:"level"`
	// 服务名称
	ServiceName string `json:"serviceName" yaml:"serviceName"`
	// 开启控制台
	// true：输出到控制台
	// 默认值：false
	OutputConsole bool `json:"outputConsole" yaml:"outputConsole"`
	// 开启文件
	// true：输出到文件
	// 默认值：false
	OutputFile bool `json:"outputFile" yaml:"outputFile"`
}

func GetLogsConfig() (*LogsConfig, error) {
	if logsConfig == nil {
		logsConfig = &LogsConfig{FileLogger: &lumberjack.Logger{}}
		if err := logsConfig.defaultValue(); err != nil {
			return nil, err
		}
	}
	return logsConfig, nil
}

// 设置日志默认值
func (l *LogsConfig) defaultValue() error {
	if l.OutputFile {
		if l.FileLogger == nil {
			l.FileLogger = &lumberjack.Logger{}
		}
		if utils.IsEmpty(l.FileLogger.Filename) {
			rootPath := getProjectRootPath()
			logsPath := filepath.Join(rootPath, logsFilePath)
			if !fileutil.Exist(logsPath) {
				if err := os.Mkdir(logsPath, os.ModePerm); err != nil {
					return err
				}
			}
			projectName := logsFilePath
			index := strings.LastIndex(rootPath, "/")
			if index != -1 {
				projectName = rootPath[index+1:]
			}
			l.FileLogger.Filename = filepath.Join(logsPath, projectName+logsFileSuffix)
		}
		if l.FileLogger.MaxAge == 0 {
			l.FileLogger.MaxAge = 7
		}
		if l.FileLogger.MaxBackups == 0 {
			l.FileLogger.MaxBackups = 10
		}
		if l.FileLogger.MaxSize == 0 {
			l.FileLogger.MaxSize = 64
		}
	}

	if utils.IsEmpty(l.Level) {
		l.Level = INFO
	}
	if utils.IsEmpty(l.ServiceName) {
		l.ServiceName = getServiceName()
	}
	return nil
}

func getServiceName() string {
	// 获取项目根路径
	rootPath := getProjectRootPath()
	index := strings.LastIndex(rootPath, "/")
	serviceName := "project"
	if index != -1 {
		serviceName = rootPath[index+1:]
	}
	return serviceName
}

func getProjectRootPath() string {
	currentPath := getCurrentPath()
	return strings.Replace(currentPath, "/core/logs", "", 1)
}

func getCurrentPath() string {
	// skip：0.表示调用者本身，获取的是当前文件名
	// skip：1.表示调用者的调用者，获取的是源头调用文件名
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
