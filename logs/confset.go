package logs

type LoggerBaseConfig struct {
	EnableFile   bool //是否输出到文件
	EnableStdout bool //是否输出到std
	LogPath      string
	LogName      string
	LogType      string
	LogLevel     string
}

type LoggerOption struct {
	LogPath              string //路径
	LogName              string //日志名
	TimeFormat           string //日志时间格式
	LogType              string //日志类型
	EnableRecordFileInfo bool   //是否记录行号
	LogLevel             Level  //日志级别
	FileInfoField        string //行号记录字段名
	EnableStdout         bool   //是否输出到stdout
	EnableFile           bool   //是否输出到文件
}

type LoggerOptionSetter struct {
	f func(option *LoggerOption)
}

func GetLogLevelByStr(s string) Level {
	l := InfoLevel
	switch s {
	case "debug":
		l = DebugLevel
	case "info":
		l = InfoLevel
	case "warning":
		l = WarnLevel
	default:
		l = InfoLevel
	}
	return l
}

func SetLogPath(p string) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.LogPath = p
	}}
}

func SetLogName(p string) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.LogName = p
	}}
}

func SetTimeFormat(p string) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.TimeFormat = p
	}}
}

func SetLogType(p string) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.LogType = p
	}}
}

func SetEnableRecordFileInfo(p bool) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.EnableRecordFileInfo = p
	}}
}

func SetLogLevel(p Level) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.LogLevel = p
	}}
}

func SetFileInfoField(p string) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.FileInfoField = p
	}}
}

func SetEnableFile(p bool) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.EnableFile = p
	}}
}

func SetEnableStdout(p bool) LoggerOptionSetter {
	return LoggerOptionSetter{func(option *LoggerOption) {
		option.EnableStdout = p
	}}
}

type LoggerConfSetter struct {
	LogOpt LoggerOption
}

func defaultOpt() LoggerOption {
	return LoggerOption{
		LogPath:              "./log",
		LogName:              "applog",
		TimeFormat:           "2006-01-02 15:04:05",
		LogType:              LOG_TYPE_TEXT,
		EnableRecordFileInfo: true,
		LogLevel:             DebugLevel,
		FileInfoField:        "call",
	}
}

func NewLoggerConfWithDefault() *LoggerConfSetter {
	return &LoggerConfSetter{LogOpt: defaultOpt()}
}

func (Self *LoggerConfSetter) Set(options ...LoggerOptionSetter) *LoggerConfSetter {
	for _, optFunc := range options {
		optFunc.f(&Self.LogOpt)
	}
	return Self
}

func (Self *LoggerConfSetter) GeConf() LoggerOption {
	return Self.LogOpt
}
