package goplogjson

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/axpira/gop/log"
)

type contextKey string

const (
	loggerContextKey contextKey = "logger"
)

var levelHook = map[log.Level]func(){
	log.PanicLevel: func() { panic("") },
	log.FatalLevel: func() { os.Exit(1) },
}

func init() {
	log.DefaultLogger = New()
}

func WithOutput(out io.Writer) log.LoggerOption {
	return log.LoggerOptionFunc(func(l log.Logger) log.Logger {
		l1 := l.(*logger)
		l1.out = out
		return l1
	})
}

func WithLevel(lv log.Level) log.LoggerOption {
	return log.LoggerOptionFunc(func(l log.Logger) log.Logger {
		l1 := l.(*logger)
		l1.level = lv
		return l1
	})
}

func newLogger() *logger {
	return &logger{
		buf:   make([]byte, 0, 500),
		level: log.InfoLevel,
		out:   os.Stdout,
	}
}

func New(opts ...log.LoggerOption) log.Logger {
	return newLogger().With(opts...)
}

type logger struct {
	buf   []byte
	level log.Level
	out   io.Writer
}

func (l *logger) clone() *logger {
	lNew := &logger{buf: make([]byte, 0, 500), level: l.level, out: l.out}
	lNew.buf = append(lNew.buf, l.buf...)
	return lNew
}

func (l *logger) Level() log.Level {
	return l.level
}

func (l *logger) HasLevel(lv log.Level) bool {
	return lv >= l.level
}

func (l *logger) NewFieldBuilder() log.FieldBuilder {
	if l.level == log.DisabledLevel {
		return emptyFieldPtr
	}
	return newField()
}

func (l *logger) With(opts ...log.LoggerOption) log.Logger {
	var logger log.Logger = l.clone()
	for _, opt := range opts {
		logger = opt.Update(logger)
	}
	return logger
}

func (l *logger) Log(lv log.Level, fieldBuilder log.FieldBuilder) {
	if fn, ok := levelHook[lv]; ok {
		defer fn()
	}
	if !l.HasLevel(lv) || fieldBuilder == emptyFieldPtr {
		if fieldBuilder != emptyFieldPtr {
			putField(fieldBuilder.(*field))
		}
		return
	}
	field := fieldBuilder.(*field)
	field.Str(LevelFieldName, LevelNameFunc(lv))
	field.buf = append(field.buf, l.buf...)
	field.send(l.out)
}

func (l *logger) Trc(f log.FieldBuilder) {
	l.Log(log.TraceLevel, f)
}

func (l *logger) Trace(msg string) {
	l.Log(log.TraceLevel, l.NewFieldBuilder().Msg(msg))
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.Log(log.TraceLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Dbg(f log.FieldBuilder) {
	l.Log(log.DebugLevel, f)
}

func (l *logger) Debug(msg string) {
	l.Log(log.DebugLevel, l.NewFieldBuilder().Msg(msg))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Log(log.DebugLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Inf(f log.FieldBuilder) {
	l.Log(log.InfoLevel, f)
}

func (l *logger) Info(msg string) {
	l.Log(log.InfoLevel, l.NewFieldBuilder().Msg(msg))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Log(log.InfoLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Wrn(f log.FieldBuilder) {
	l.Log(log.WarnLevel, f)
}

func (l *logger) Warn(msg string) {
	l.Log(log.WarnLevel, l.NewFieldBuilder().Msg(msg))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Log(log.WarnLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Err(f log.FieldBuilder) {
	l.Log(log.ErrorLevel, f)
}

func (l *logger) Error(msg string, err error) {
	l.Log(log.ErrorLevel, l.NewFieldBuilder().Msg(msg).Err(err))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Log(log.ErrorLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Ftl(f log.FieldBuilder) {
	l.Log(log.FatalLevel, f)
}

func (l *logger) Fatal(msg string) {
	l.Log(log.FatalLevel, l.NewFieldBuilder().Msg(msg))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Log(log.FatalLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Pnc(f log.FieldBuilder) {
	l.Log(log.PanicLevel, f)
}

func (l *logger) Panic(msg string) {
	l.Log(log.PanicLevel, l.NewFieldBuilder().Msg(msg))
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.Log(log.PanicLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Print(args ...interface{}) {
	l.Log(log.InfoLevel, l.NewFieldBuilder().Msg(fmt.Sprint(args...)))
}

func (l *logger) Println(args ...interface{}) {
	l.Log(log.InfoLevel, l.NewFieldBuilder().Msg(fmt.Sprint(args...)))
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.Log(log.InfoLevel, l.NewFieldBuilder().Msgf(format, args...))
}

func (l *logger) Write(msg []byte) (int, error) {
	l.Log(log.InfoLevel, l.NewFieldBuilder().Msg(string(msg)))
	return 0, nil
}

func (l *logger) FromCtx(ctx context.Context) log.Logger {
	lg := ctx.Value(loggerContextKey)
	if lg == nil {
		return l
	}
	return lg.(*logger)
}

func (l *logger) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}
