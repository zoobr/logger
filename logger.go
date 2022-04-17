package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerModeDev     = "dev"          // development logger mode
	loggerModeProd    = "prod"         // production logger mode
	loggerModeTesting = "testing"      // logger mode for tests
	defaultLoggerMode = loggerModeProd // by default apply prod mode
)

// public main packet functions
var (

	// JSONEncoder JSON log format
	JSONEncoder = 0
	// ConsoleEncoder Console log format
	ConsoleEncoder = 1

	// methods wihout template
	// LogMethodDebug debug mode
	LogMethodDebug = 1
	// LogMethodInfo info mode
	LogMethodInfo = 2
	// LogMethodWarn warn mode
	LogMethodWarn = 3
	// LogMethodError error mode
	LogMethodError = 4
	// LogMethodPanic panic mode
	LogMethodPanic = 5

	// methods with templates
	// LogMethodDebugf debugf mode
	LogMethodDebugf = 6
	// LogMethodDebugw debugw mode
	LogMethodDebugw = 7
	// LogMethodInfof infof mode
	LogMethodInfof = 8
	// LogMethodWarnf warnf mode
	LogMethodWarnf = 9
	// LogMethodErrorf errorf mode
	LogMethodErrorf = 10
	// LogMethodPanicf panicf mode
	LogMethodPanicf = 11

	// public packet function for use logger package without create separated logger instance
	// Debug
	Debug  func(...interface{})
	Debugf func(string, ...interface{})
	Debugw func(string, ...interface{})
	Info   func(...interface{})
	Infof  func(string, ...interface{})
	Warn   func(...interface{})
	Warnf  func(string, ...interface{})
	Error  func(...interface{})
	Errorf func(string, ...interface{})
	Panic  func(...interface{})
	Panicf func(string, ...interface{})
)

// internal vars for inner logic
var (
	emptyTemplate = "<>" // this template is used if for some reason the template for logging is not passed
	logger        *Sugar
	// default config
	defaultConfig = Config{
		LoggerMode:  "prod",
		EncoderType: 1,
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.RFC3339TimeEncoder,
	}
)

// Config struct for init logger configaration
type Config struct {
	LoggerType  string
	LoggerMode  string
	EncoderType int
	EncodeLevel zapcore.LevelEncoder
	EncodeTime  zapcore.TimeEncoder
}

// Logger structure for extending the functionality of a standard logger
type Logger struct {
	*zap.Logger
	applyTemplate func(*string)
}

// Log basic method providing logging. All logging methods work through this function
func (logger *Logger) Log(method int, template *string, fields ...zapcore.Field) {
	if logger.applyTemplate != nil {
		logger.applyTemplate(template)
	}
	log := logger.Logger
	if template == nil {
		log.Warn("use logger method without template", fields...)
		template = &emptyTemplate
	}
	switch method {
	case LogMethodDebug:
		log.Debug(*template, fields...)
	case LogMethodInfo:
		log.Info(*template, fields...)
	case LogMethodWarn:
		log.Warn(*template, fields...)
	case LogMethodError:
		log.Error(*template, fields...)
	case LogMethodPanic:
		log.Panic(*template, fields...)
	}
}

// Debug loggin debug messages
func (logger *Logger) Debug(template string, args ...zapcore.Field) {
	logger.Log(LogMethodDebug, &template)
}

// Error loggin errors messages
func (logger *Logger) Error(template string, args ...zapcore.Field) {
	logger.Log(LogMethodError, &template)
}

// Warn loggin errors messages
func (logger *Logger) Warn(template string, args ...zapcore.Field) {
	logger.Log(LogMethodWarn, &template)
}

// Panic loggin messages and run panic
func (logger *Logger) Panic(template string, args ...zapcore.Field) {
	logger.Log(LogMethodPanic, &template)
}

// Sync flushes any buffered log entries
func (logger *Logger) Sync() {
	logger.Logger.Sync()
}

// private constructor for create logger
func createLogger(config *Config) *zap.Logger {
	core := prepareConfig(config)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))
}

// NewLogger constructor for create logger
func NewLogger(config *Config) *zap.Logger {
	return createLogger(config)
}

// SetApplyTemplate Sets a callback that is called every time a message is generated for logging.
// A link to the generated message template is passed to the function.
// Thus, you can add, change the line of logs in a template way at your discretion
func (logger *Logger) SetApplyTemplate(cb func(*string)) {
	logger.applyTemplate = cb
}

// Init prepare logger structure
func Init(config *Config) {
	logger = createSugaredLogger(config)
}

func init() {
	logger = createSugaredLogger(nil)
	Debug = logger.Debug
	Debugf = logger.Debugf
	Debugw = logger.Debugw
	Error = logger.Error
	Errorf = logger.Errorf
	Info = logger.Info
	Infof = logger.Infof
	Warn = logger.Warn
	Warnf = logger.Warnf
	Panic = logger.Panic
	Panicf = logger.Panicf
}

func prepareConfig(config *Config) zapcore.Core {
	// prepare config
	if config == nil {
		config = &defaultConfig
	} else {
		// check config params and set defaults is empty
		if config.EncodeLevel == nil {
			config.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		if config.EncodeTime == nil {
			config.EncodeTime = zapcore.RFC3339TimeEncoder
		}
	}

	// prepare logger mode
	var configEncoder zapcore.EncoderConfig
	logLevel := zapcore.DebugLevel
	loggerMode := config.LoggerMode
	if loggerMode == loggerModeProd { // logger for development mode
		configEncoder = zap.NewProductionEncoderConfig()
		logLevel = zapcore.InfoLevel
	} else {
		if loggerMode != loggerModeDev && len(loggerMode) > 0 {
			fmt.Printf("wrong logger mode: %s, will use dev logger", loggerMode)
		} else if len(loggerMode) == 0 {
			fmt.Printf("logger mode is empty, will use dev logger")
		}
		configEncoder = zap.NewDevelopmentEncoderConfig()
	}

	configEncoder.EncodeLevel = config.EncodeLevel
	configEncoder.EncodeTime = config.EncodeTime

	// prepare encoder
	var newEncoder zapcore.Encoder
	switch config.EncoderType {
	case ConsoleEncoder:
		newEncoder = zapcore.NewConsoleEncoder(configEncoder)
	default:
		newEncoder = zapcore.NewJSONEncoder(configEncoder)
	}
	core := zapcore.NewCore(newEncoder, os.Stdout, logLevel)
	return core
}
