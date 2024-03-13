package log

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/utils"
	"time"
)

var Zap = new(_zap)

type _zap struct{}

func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(global.CONF.Zap.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
}

func (z *_zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.CONF.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    global.CONF.Zap.ZapEncodeLevel(),
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func (z *_zap) GetEncoder() zapcore.Encoder {
	if global.CONF.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

func (z *_zap) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := FileRotatelogs.GetWriteSyncer(l.String())     // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}
	return zapcore.NewCore(z.GetEncoder(), writer, level)
}

func (z *_zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := global.CONF.Zap.TransportLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
	}
	return cores
}

func (z *_zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

// GetWriteSyncer get zapcore.WriteSyncer
func (r *fileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(global.CONF.Zap.Director, "%Y-%m-%d", level+".log"),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(global.CONF.Zap.MaxAge)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if global.CONF.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

//NewZap get zap.Logger
func NewZap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.CONF.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", global.CONF.Zap.Director)
		_ = os.MkdirAll(global.CONF.Zap.Director, os.ModePerm)
	}
	cores := Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.CONF.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}