package logs

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
)

var (
	LOG_TYPE_TEXT = "text"
	LOG_TYPE_JSON = "json"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func newLogger(options LoggerOption) (*logrus.Logger, error) {
	var (
		err          error
		logrusLogger *logrus.Logger
	)
	//创建文件夹
	if err = makeLogDir(options.LogPath); err != nil {
		return logrusLogger, err
	}
	logrusLogger = logrus.New()
	logrusLogger.SetLevel(logrus.Level(options.LogLevel))
	logrusLogger.Out = os.Stdout
	switch options.LogType {
	case LOG_TYPE_JSON:
		format := &logrus.JSONFormatter{
			TimestampFormat: options.TimeFormat,
		}
		logrusLogger.Formatter = format
	default:
		logrusLogger.Formatter = &logrus.TextFormatter{
			TimestampFormat: options.TimeFormat,
		}
	}
	return logrusLogger, nil
}

type logger struct {
	write                io.Writer
	logrus               *logrus.Logger
	enableRecordFileInfo bool
	fileInfoField        string
	env                  string
}

func New(loggerOpt LoggerOption) (*logger, error) {
	var (
		err          error
		logrusLogger *logrus.Logger
	)
	if logrusLogger, err = newLogger(loggerOpt); err != nil {
		return nil, err
	}
	//日志轮转
	absPath, err := filepath.Abs(loggerOpt.LogPath)
	filePath := path.Join(absPath, loggerOpt.LogName)
	src, err := os.OpenFile(filePath, os.O_CREATE|os.O_SYNC|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("日志文件打开失败:", err)
		os.Exit(-1)
	}
	wrs := make([]io.Writer, 0)
	if loggerOpt.EnableStdout {
		wrs = append(wrs, colorable.NewColorableStdout())
	}
	if loggerOpt.EnableFile {
		wrs = append(wrs, src)
	}
	mw := io.MultiWriter(wrs...)
	logrusLogger.SetOutput(mw)
	logger := &logger{
		write:                mw,
		logrus:               logrusLogger,
		enableRecordFileInfo: loggerOpt.EnableRecordFileInfo,
		fileInfoField:        loggerOpt.FileInfoField,
	}
	return logger, nil
}

type D map[string]interface{}

func (l *logger) GetLogger() *logrus.Logger {
	return l.logrus
}

func (l *logger) GetWriter() io.Writer {
	return l.write
}

//这些接口可以通过logger自行进行符合业务逻辑的二次开发
func (l *logger) Debug(dataFields D, message string) {
	if dataFields == nil {
		dataFields = D{}
	}
	if l.enableRecordFileInfo {
		dataFields[l.fileInfoField] = fileInfo(2)
	}
	l.logrus.WithFields(logrus.Fields(dataFields)).Debug(message)
}

func (l *logger) Info(dataFields D, message string) {
	if dataFields == nil {
		dataFields = D{}
	}
	if l.enableRecordFileInfo {
		dataFields[l.fileInfoField] = fileInfo(2)
	}
	l.logrus.WithFields(logrus.Fields(dataFields)).Info(message)
}

func (l *logger) Warn(dataFields D, message string) {
	if dataFields == nil {
		dataFields = D{}
	}
	if l.enableRecordFileInfo {
		dataFields[l.fileInfoField] = fileInfo(2)
	}
	l.logrus.WithFields(logrus.Fields(dataFields)).Warn(message)
}

func (l *logger) Error(dataFields D, message string) {
	if dataFields == nil {
		dataFields = D{}
	}
	if l.enableRecordFileInfo {
		dataFields[l.fileInfoField] = fileInfo(2)
	}
	l.logrus.WithFields(logrus.Fields(dataFields)).Error(message)
}

func (l *logger) Fatal(dataFields D, message string) {
	if dataFields == nil {
		dataFields = D{}
	}
	if l.enableRecordFileInfo {
		dataFields[l.fileInfoField] = fileInfo(2)
	}
	l.logrus.WithFields(logrus.Fields(dataFields)).Fatal(message)
}

func (l *logger) Panic(dataFields D, message string) {
	if dataFields == nil {
		dataFields = D{}
	}
	if l.enableRecordFileInfo {
		dataFields[l.fileInfoField] = fileInfo(2)
	}
	l.logrus.WithFields(logrus.Fields(dataFields)).Panic(message)
}

var Logger *logger

//初始化单例
func InitLogger(loggerOpt LoggerOption, env string) {
	l, err := New(loggerOpt)
	if err != nil {
		Logger = &logger{}
		return
	}
	l.env = env
	Logger = l
}

//输出debug
func Debug(args interface{}) {
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Debug(args)
}

//输出info
func Info(args interface{}) {
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Info(args)
}

//输出wran
func Warn(args interface{}) {
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Warn(args)
}

//输出错误
func Error(args interface{}, err error) {
	if Logger.env == "product" {
		raven.CaptureErrorAndWait(errors.New(fmt.Sprintf("args:%v error:%v", args, err)), nil)
	}
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	dataFields["error"] = err
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Error(args)
}

//格式化输出
func DebugF(format string, args ...interface{}) {
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Debugf(format, args)
}

//输出info
func InfoF(format string, args ...interface{}) {
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Infof(format, args)
}

//输出wran
func WarnF(format string, args ...interface{}) {
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Warnf(format, args)
}

//输出错误
func ErrorF(err error, format string, args ...interface{}) {
	if Logger.env == "product" { //生产模式下会到sentry
		raven.CaptureErrorAndWait(errors.New(fmt.Sprintf(format, args)), nil)
	}
	dataFields := D{}
	if Logger.enableRecordFileInfo {
		dataFields[Logger.fileInfoField] = fileInfo(2)
	}
	dataFields["error"] = err
	Logger.logrus.WithFields(logrus.Fields(dataFields)).Errorf(format, args)
}
