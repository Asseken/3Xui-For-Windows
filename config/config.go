package config

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

//go:embed version
var version string

//go:embed name
var name string

type LogLevel string

const (
	Debug  LogLevel = "debug"
	Info   LogLevel = "info"
	Notice LogLevel = "notice"
	Warn   LogLevel = "warn"
	Error  LogLevel = "error"
)

func GetVersion() string {
	return strings.TrimSpace(version)
}

func GetName() string {
	return strings.TrimSpace(name)
}

func GetLogLevel() LogLevel {
	if IsDebug() {
		return Debug
	}
	logLevel := os.Getenv("XUI_LOG_LEVEL")
	if logLevel == "" {
		return Info
	}
	return LogLevel(logLevel)
}

func IsDebug() bool {
	return os.Getenv("XUI_DEBUG") == "true"
}

// xui工作目录
func GetBinFolderPath() string {
	binFolderPath := os.Getenv("XUI_BIN_FOLDER")
	if binFolderPath == "" {
		binFolderPath = "etc"
	}
	return binFolderPath
}

// x-ui HTML
func GetHtmlPath() string {
	XrayHtmlPath := os.Getenv("XUI_HTML_FOLDER")
	if XrayHtmlPath == "" {
		XrayHtmlPath = "etc/html"
	}
	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(XrayHtmlPath); os.IsNotExist(err) {
		err := os.MkdirAll(XrayHtmlPath, os.ModePerm)
		if err != nil {
			// 处理创建目录失败的错误
			panic(err)
		}
	}
	return XrayHtmlPath
}

// -------new xray file------------------
func GetXrayFolderPath() string {
	XrayFolderPath := os.Getenv("XUI_BIN_FOLDER")
	if XrayFolderPath == "" {
		XrayFolderPath = "etc/xray"
	}
	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(XrayFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(XrayFolderPath, os.ModePerm)
		if err != nil {
			// 处理创建目录失败的错误
			panic(err)
		}
	}
	return XrayFolderPath
}

func GetDBFolderPath() string {
	dbFolderPath := os.Getenv("XUI_DB_FOLDER")
	if dbFolderPath == "" {
		dbFolderPath = "etc/x-ui"
	}
	return dbFolderPath
}

func GetDBPath() string {
	return fmt.Sprintf("%s/%s.db", GetDBFolderPath(), GetName())
}

func GetLogFolder() string {
	logFolderPath := os.Getenv("XUI_LOG_FOLDER")
	if logFolderPath == "" {
		logFolderPath = "etc/log"
	}
	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		err := os.MkdirAll(logFolderPath, os.ModePerm)
		if err != nil {
			// 处理创建目录失败的错误
			panic(err)
		}
	}
	return logFolderPath
}
