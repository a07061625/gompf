/**
 * 日志
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 9:22
 */
package mplog

import (
    "encoding/json"
    "log"
    "os"
    "sync"
    "time"

    "github.com/robfig/cron"
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

const (
    LogTypeDebug  = "DEBUG"
    LogTypeWarn   = "WARN"
    LogTypeInfo   = "INFO"
    LogTypeError  = "ERROR"
    LogTypeDPanic = "DPANIC"
    LogTypePanic  = "PANIC"
    LogTypeFatal  = "FATAL"
)

type LogField struct {
    Key string
    Val interface{}
}

type logWriter struct {
    file *os.File
}

func (f *logWriter) Write(p []byte) (n int, err error) {
    return f.file.Write(p)
}

func (f *logWriter) Close() error {
    return f.file.Close()
}

type logDaily struct {
    logger       *zap.Logger
    loggerFields map[string]interface{}
    accessFile   *logWriter
    errorFile    *logWriter
    cron         *cron.Cron
    logPrefix    string
    logDir       string
    logAccess    string
    logError     string
    logSuffix    string
}

func (ld *logDaily) createLogger() {
    highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level >= zap.ErrorLevel
    })
    lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return (level < zap.ErrorLevel) && (level >= zap.DebugLevel)
    })

    prodEncoder := zap.NewProductionEncoderConfig()
    prodEncoder.TimeKey = ""
    prodEncoder.LevelKey = ""
    prodEncoder.NameKey = ""
    prodEncoder.CallerKey = ""
    prodEncoder.StacktraceKey = "S"
    prodEncoder.MessageKey = "M"
    prodEncoder.LineEnding = zapcore.DefaultLineEnding
    prodEncoder.EncodeDuration = zapcore.StringDurationEncoder
    prodEncoder.EncodeName = zapcore.FullNameEncoder

    // 开启开发模式,调用跟踪
    caller := zap.AddCaller()
    // 开启堆栈跟踪
    stacktrace := zap.AddStacktrace(zapcore.ErrorLevel)
    // 开启文件及行号
    development := zap.Development()
    highCore := zapcore.NewCore(zapcore.NewConsoleEncoder(prodEncoder), zapcore.AddSync(ld.errorFile), highPriority)
    lowCore := zapcore.NewCore(zapcore.NewConsoleEncoder(prodEncoder), zapcore.AddSync(ld.accessFile), lowPriority)
    ld.logger = zap.New(zapcore.NewTee(highCore, lowCore), caller, stacktrace, development)
}

func (ld *logDaily) getLogName(level string) string {
    fileName := ""
    if level == "error" {
        fileName = ld.logError
    } else {
        fileName = ld.logAccess
    }
    fileName += "_" + time.Now().Format("20060102") + ld.logSuffix
    logName := ld.logDir + "/" + fileName

    fileInfo, err := os.Stat(logName)
    if os.IsNotExist(err) {
        f, err := os.Create(logName)
        if err != nil {
            log.Fatalln("create log file error:" + err.Error())
        }
        defer f.Close()
    } else if err != nil {
        log.Fatalln("log file error:" + err.Error())
    } else if fileInfo.IsDir() {
        log.Fatalln("log file is dir")
    }

    return logName
}

func (ld *logDaily) ChangeAccessLog() {
    if ld.accessFile.file != nil {
        ld.accessFile.file.Close()
    }

    infoOutput, err := os.OpenFile(ld.getLogName("info"), os.O_RDWR|os.O_APPEND, 0666)
    if err == nil {
        ld.accessFile.file = infoOutput
    } else {
        log.Fatalln("access log error:" + err.Error())
    }
}

func (ld *logDaily) ChangeErrorLog() {
    if ld.errorFile.file != nil {
        ld.errorFile.file.Close()
    }

    errOutput, err := os.OpenFile(ld.getLogName("error"), os.O_RDWR|os.O_APPEND, 0666)
    if err == nil {
        ld.errorFile.file = errOutput
    } else {
        log.Fatalln("error log error:" + err.Error())
    }
}

func (ld *logDaily) log(level, msg string, fields ...LogField) {
    fieldList := make(map[string]interface{})
    if len(ld.loggerFields) > 0 {
        for k, v := range ld.loggerFields {
            fieldList[k] = v
        }
    }
    if len(fields) > 0 {
        for _, eField := range fields {
            fieldList[eField.Key] = eField.Val
        }
    }

    fieldJson, _ := json.Marshal(fieldList)
    logMsg := time.Now().Format("2006-01-02 03:04:05.000") + " | " + level + " | " + os.Getenv("MP_REQ_ID") + ld.logPrefix
    logMsg += " | " + string(fieldJson) + "\n" + msg
    switch level {
    case LogTypeInfo:
        ld.logger.Info(logMsg)
    case LogTypeError:
        ld.logger.Error(logMsg)
    case LogTypeFatal:
        ld.logger.Fatal(logMsg)
    case LogTypeWarn:
        ld.logger.Warn(logMsg)
    case LogTypeDebug:
        ld.logger.Debug(logMsg)
    case LogTypeDPanic:
        ld.logger.DPanic(logMsg)
    case LogTypePanic:
        ld.logger.Panic(logMsg)
    default:
        ld.logger.Info(logMsg)
    }
}

var (
    once sync.Once
    ins  *logDaily
)

func init() {
    ins = &logDaily{}
    ins.loggerFields = make(map[string]interface{})
    ins.accessFile = &logWriter{}
    ins.errorFile = &logWriter{}
    ins.cron = cron.New()
}

func Load(conf *viper.Viper, fields map[string]interface{}, extend map[string]interface{}) {
    once.Do(func() {
        ins.logPrefix = " | " + extend["server_host"].(string)
        ins.logPrefix += " | " + extend["server_port"].(string)
        ins.logPrefix += " | " + extend["env_type"].(string) + extend["project_tag"].(string)
        ins.logPrefix += " | " + extend["project_module"].(string)

        ins.logDir = extend["log_dir"].(string)
        err := os.MkdirAll(ins.logDir, os.ModePerm)
        if err != nil {
            log.Fatalln("log dir create fail:" + err.Error())
        }
        confPrefix := extend["conf_prefix"].(string)
        ins.logAccess = conf.GetString(confPrefix + "access")
        ins.logError = conf.GetString(confPrefix + "error")
        ins.logSuffix = conf.GetString(confPrefix + "suffix")

        ins.cron.AddFunc(conf.GetString(confPrefix+"cron.access"), ins.ChangeAccessLog)
        ins.cron.AddFunc(conf.GetString(confPrefix+"cron.error"), ins.ChangeErrorLog)
        ins.cron.Start()
        ins.ChangeAccessLog()
        ins.ChangeErrorLog()

        ins.loggerFields = conf.GetStringMap(confPrefix + "fields")
        if len(fields) > 0 {
            for k, v := range fields {
                ins.loggerFields[k] = v
            }
        }
        ins.createLogger()
    })
}

func LogDebug(msg string, fields ...LogField) {
    ins.log(LogTypeDebug, msg, fields...)
}

func LogWarn(msg string, fields ...LogField) {
    ins.log(LogTypeWarn, msg, fields...)
}

func LogInfo(msg string, fields ...LogField) {
    ins.log(LogTypeInfo, msg, fields...)
}

func LogError(msg string, fields ...LogField) {
    ins.log(LogTypeError, msg, fields...)
}

func LogDPanic(msg string, fields ...LogField) {
    ins.log(LogTypeDPanic, msg, fields...)
}

func LogPanic(msg string, fields ...LogField) {
    ins.log(LogTypePanic, msg, fields...)
}

func LogFatal(msg string, fields ...LogField) {
    ins.log(LogTypeFatal, msg, fields...)
}
