package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	WithField(string, interface{}) Logger
	SetLevel(int)
}

func New() Logger {
	return &l{*logrus.New()}
}

type l struct {
	logrus.Logger
}

type e struct {
	logrus.Entry
}

func (l *l) WithField(key string, value interface{}) Logger {
	entry := l.Logger.WithField(key, value)
	return &e{*entry}
}

func (l *l) SetLevel(level int) {
	l.Logger.SetLevel(logrus.Level(level))
}

func (en *e) WithField(key string, value interface{}) Logger {
	entry := en.Entry.Logger.WithField(key, value)
	return &e{*entry}
}

func (e *e) SetLevel(level int) {
	e.Logger.SetLevel(logrus.Level(level))
}
