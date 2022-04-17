package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Sugar structure for extending the functionality of a standard sugared logger
type Sugar struct {
	*zap.SugaredLogger
	applyTemplate func(*string)
}

// Log basic method providing logging. All logging methods work through this function
func (logger *Sugar) Log(method int, template *string, args ...interface{}) {
	if logger.applyTemplate != nil {
		logger.applyTemplate(template)
	}
	log := logger.SugaredLogger
	if method > LogMethodPanic && template == nil {
		log.Warnf("use logger method without template", args...)
		template = &emptyTemplate
	}
	switch method {
	case LogMethodDebug:
		log.Debug(args...)
	case LogMethodDebugf:
		log.Debugf(*template, args...)
	case LogMethodDebugw:
		log.Debugw(*template, args...)
	case LogMethodInfo:
		log.Info(args...)
	case LogMethodInfof:
		log.Infof(*template, args...)
	case LogMethodWarn:
		log.Warn(args...)
	case LogMethodWarnf:
		log.Warnf(*template, args...)
	case LogMethodError:
		log.Error(args...)
	case LogMethodErrorf:
		log.Errorf(*template, args...)
	case LogMethodPanic:
		log.Panic(args...)
	case LogMethodPanicf:
		log.Panicf(*template, args...)
	}
}

// Debug loggin debug messages
func (logger *Sugar) Debug(args ...interface{}) {
	logger.Log(LogMethodDebug, nil, args...)
}

// Debugf loggin debug messages whit special template
// Debugf uses fmt.Sprintf to log a templated message
func (logger *Sugar) Debugf(template string, args ...interface{}) {
	logger.Log(LogMethodDebugf, &template, args...)
}

// Debugw loggin debug messages
// Debugw logs a message with some additional context. The variadic key-value pairs are treated as they are in With
func (logger *Sugar) Debugw(template string, keyAndValues ...interface{}) {
	logger.Log(LogMethodDebugw, &template, keyAndValues...)
}

// Warn loggin errors messages
func (logger *Sugar) Warn(args ...interface{}) {
	logger.Log(LogMethodWarn, nil, args...)
}

// Warnf loggin errors messages
// Warnf uses fmt.Sprintf to log a templated message.
func (logger *Sugar) Warnf(template string, args ...interface{}) {
	logger.Log(LogMethodWarnf, &template, args...)
}

// Error loggin errors messages
func (logger *Sugar) Error(args ...interface{}) {
	logger.Log(LogMethodError, nil, args...)
}

// Errorf loggin errors messages
// Errorf uses fmt.Sprintf to log a templated message
func (logger *Sugar) Errorf(template string, args ...interface{}) {
	logger.Log(LogMethodErrorf, &template, args...)
}

// Panic loggin messages and run panic
func (logger *Sugar) Panic(args ...interface{}) {
	logger.Log(LogMethodPanic, nil, args...)
}

// Panicf loggin messages and run panic
// Panicf uses fmt.Sprintf to log a templated message
func (logger *Sugar) Panicf(template string, args ...interface{}) {
	logger.Log(LogMethodPanicf, &template, args...)
}

// Sync flushes any buffered log entries
func (logger *Sugar) Sync() {
	logger.SugaredLogger.Sync()
}

// private constructor for create sugared logger
func createSugaredLogger(config *Config) *Sugar {
	core := prepareConfig(config)
	return &Sugar{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel)).Sugar(), nil}
}

// NewSugaredLogger constructor for create sugared logger
func NewSugaredLogger(config *Config) *Sugar {
	return createSugaredLogger(config)
}

// SetApplyTemplate Sets a callback that is called every time a message is generated for logging.
// A link to the generated message template is passed to the function.
// Thus, you can add, change the line of logs in a template way at your discretion
func (logger *Sugar) SetApplyTemplate(cb func(*string)) {
	logger.applyTemplate = cb
}
