/**
 * 日志
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 9:22
 */
package mpf

import (
    "log"
    "os"
    "time"

    "github.com/robfig/cron"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

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
    Logger     *zap.Logger
    accessFile *logWriter
    errorFile  *logWriter
    cron       *cron.Cron
    logDir     string
    logAccess  string
    logError   string
    logSuffix  string
}

func logEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func logEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
}

func (ld *logDaily) createLogger(initFields map[string]string) {
    highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return level >= zap.ErrorLevel
    })
    lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
        return (level < zap.ErrorLevel) && (level >= zap.DebugLevel)
    })

    prodEncoder := zap.NewProductionEncoderConfig()
    prodEncoder.TimeKey = "T"
    prodEncoder.LevelKey = "L"
    prodEncoder.NameKey = "N"
    prodEncoder.CallerKey = "C"
    prodEncoder.MessageKey = "M"
    prodEncoder.StacktraceKey = "S"
    prodEncoder.LineEnding = zapcore.DefaultLineEnding
    prodEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
    prodEncoder.EncodeTime = logEncodeTime
    prodEncoder.EncodeDuration = zapcore.StringDurationEncoder
    prodEncoder.EncodeCaller = logEncodeCaller
    prodEncoder.EncodeName = zapcore.FullNameEncoder

    // 开启开发模式,调用跟踪
    caller := zap.AddCaller()
    // 开启堆栈跟踪
    stacktrace := zap.AddStacktrace(zapcore.ErrorLevel)
    // 开启文件及行号
    development := zap.Development()
    // 设置初始化字段
    arr := make([]zap.Field, 0)
    for key, val := range initFields {
        arr = append(arr, zap.String(key, val))
    }
    fields := zap.Fields(arr...)
    highCore := zapcore.NewCore(zapcore.NewConsoleEncoder(prodEncoder), zapcore.AddSync(ld.errorFile), highPriority)
    lowCore := zapcore.NewCore(zapcore.NewConsoleEncoder(prodEncoder), zapcore.AddSync(ld.accessFile), lowPriority)
    ld.Logger = zap.New(zapcore.NewTee(highCore, lowCore), caller, stacktrace, development, fields)
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

func (ld *logDaily) SetCron(cron *cron.Cron) {
    ld.cron = cron
}

func (ld *logDaily) SetLogAccess(logAccess string) {
    ld.logAccess = logAccess
}

func (ld *logDaily) SetLogError(logError string) {
    ld.logError = logError
}

func (ld *logDaily) SetLogSuffix(logSuffix string) {
    ld.logSuffix = logSuffix
}

var (
    insLog *logDaily
)

func init() {
    insLog = &logDaily{nil, &logWriter{}, &logWriter{}, nil, "", "", "", ""}
}

func NewLogger() *zap.Logger {
    return insLog.Logger
}
