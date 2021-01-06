package logger

import (
	"github.com/mitchellh/go-homedir"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path"
)

var logger *Logger

type Logger struct {
	*zap.SugaredLogger
}

func init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	log := zap.New(core,zap.AddCaller())
	logger = &Logger{log.Sugar()}
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	dir, err := homedir.Expand("~/.phoenix")
	if err != nil {
		log.Fatal(err)
	}
	destination := path.Join(dir, "test.log")
	os.Create(destination)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   destination,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func (self *Logger) Write(b []byte) (n int, err error) {
	self.Info(string(b))
	return len(b), nil
}

func GetLogger() *Logger {
	return logger
}

func SetLogger(newLogger *Logger) {
	defer logger.Sync()
	logger = newLogger
}


func LoggerWriter() *Logger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	log := zap.New(core/*,zap.AddCaller()*/)
	logger = &Logger{log.Sugar()}
	return logger
}
