package logs

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogrusWriter struct {
	logger *logrus.Logger
	ctx    context.Context
}

func NewLogrusWriter(opts ...func(logger *logrus.Logger)) logx.Writer {
	logger := logrus.New()
	for _, opt := range opts {
		opt(logger)
	}

	return &LogrusWriter{
		logger: logger,
	}
}

func (w *LogrusWriter) Alert(v interface{}) {
	w.logger.Error(v)
}

func (w *LogrusWriter) Close() error {
	w.logger.Exit(0)
	return nil
}

func (w *LogrusWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Debug(v)
}

func (w *LogrusWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Error(v)
}

func (w *LogrusWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Info(v)
}

func (w *LogrusWriter) Severe(v interface{}) {
	w.logger.Fatal(v)
}

func (w *LogrusWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Warn(v)
}

func (w *LogrusWriter) Stack(v interface{}) {
	w.logger.Error(v)
}

func (w *LogrusWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.WithFields(toLogrusFields(fields...)).Info(v)
}

func toLogrusFields(fields ...logx.LogField) logrus.Fields {
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}

func Error(ctx context.Context, v ...any) {
	getLogger(ctx).Error(v...)
}

func Info(ctx context.Context, v ...any) {
	getLogger(ctx).Info(v...)
}

func getLogger(ctx context.Context) logx.Logger {
	return logx.WithContext(ctx).WithCallerSkip(1)
}

// func fieldCtx(ctx context.Context) logx.LogField {
// 	// return logx.Field(key, value)
// }
