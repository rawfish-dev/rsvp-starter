package logger

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/interfaces"

	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

type Logger struct {
	contextID  string
	loggerImpl *logrus.Logger
}

func NewLogger() interfaces.Logger {
	return &Logger{
		contextID:  uuid.NewV4().String(),
		loggerImpl: logrus.New(),
	}
}

func NewLoggerWithContext(contextID string) interfaces.Logger {
	return &Logger{
		contextID:  contextID,
		loggerImpl: logrus.New(),
	}
}

func (l *Logger) generateContextMessage() string {
	return fmt.Sprintf("[%v] :: ", l.contextID)
}

func (l *Logger) Info(args ...interface{}) {
	l.loggerImpl.Infof(l.generateContextMessage()+"%v", args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.loggerImpl.Infof(l.generateContextMessage()+format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.loggerImpl.Warnf(l.generateContextMessage()+"%v", args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.loggerImpl.Warnf(l.generateContextMessage()+format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.loggerImpl.Errorf(l.generateContextMessage()+"%v", args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.loggerImpl.Errorf(l.generateContextMessage()+format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.loggerImpl.Fatalf(l.generateContextMessage()+"%v", args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.loggerImpl.Fatalf(l.generateContextMessage()+format, args...)
}
