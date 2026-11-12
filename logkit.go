// Package logkit 在 go-zero logx 之上提供完整日志级别，并新增 Warn。
//
//	Debug  调试
//	Info   普通
//	Warn   预期内错误，不告警，按 biz_level=warn 检索
//	Error  需要及时处理的错误，按 level=error 接入告警
//	Severe 致命
//	Slow   慢日志
//	Stat   统计
package logkit

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// LevelKey 是 Warn 日志的等级字段，检索/告警系统按它筛预期内错误。
	LevelKey = "biz_level"
	// WarnLevel 是 Warn 日志的等级值。
	WarnLevel = "warn"
)

// Default 是包级函数使用的 logger。
var Default = New()

// Fields 把 kv 键值对转成 logx.LogField；error 值转成字符串，奇数个参数补 nil。
func Fields(kv ...any) []logx.LogField {
	if len(kv)%2 != 0 {
		kv = append(kv, nil)
	}

	fields := make([]logx.LogField, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		fields = append(fields, logx.Field(fmt.Sprint(kv[i]), fieldValue(kv[i+1])))
	}
	return fields
}

func fieldValue(v any) any {
	if err, ok := v.(error); ok && err != nil {
		return err.Error()
	}
	return v
}

// Logger 携带一组公共字段，可注入到其他组件。
type Logger struct {
	fields []logx.LogField
}

// New 创建带公共字段的 logger。
func New(kv ...any) *Logger {
	return &Logger{fields: Fields(kv...)}
}

// With 返回追加公共字段的新 logger，原 logger 不变。
func (l *Logger) With(kv ...any) *Logger {
	fields := make([]logx.LogField, 0, len(l.fields)+len(kv)/2)
	fields = append(fields, l.fields...)
	fields = append(fields, Fields(kv...)...)
	return &Logger{fields: fields}
}

// Debug 记录调试日志，Level=info 时默认不输出。
func (l *Logger) Debug(ctx context.Context, msg string, kv ...any) {
	logx.WithContext(ctx).Debugw(msg, l.merge(kv...)...)
}

// Info 记录普通日志。
func (l *Logger) Info(ctx context.Context, msg string, kv ...any) {
	logx.WithContext(ctx).Infow(msg, l.merge(kv...)...)
}

// Warn 记录预期内错误，自动附加 biz_level=warn。
func (l *Logger) Warn(ctx context.Context, msg string, kv ...any) {
	logx.WithContext(ctx).Infow(msg, l.warnFields(kv...)...)
}

// Error 记录需要及时处理的错误，写入 error 级别。
func (l *Logger) Error(ctx context.Context, msg string, kv ...any) {
	logx.WithContext(ctx).Errorw(msg, l.merge(kv...)...)
}

// Severe 记录致命日志。logx 的 Severe 不支持字段，字段会拼进消息。
func (l *Logger) Severe(_ context.Context, msg string, kv ...any) {
	logx.Severe(render(msg, l.merge(kv...)))
}

// Slow 记录慢日志。
func (l *Logger) Slow(ctx context.Context, msg string, kv ...any) {
	logx.WithContext(ctx).Sloww(msg, l.merge(kv...)...)
}

// Stat 记录统计日志。logx 的 Stat 不支持字段，字段会拼进消息。
func (l *Logger) Stat(_ context.Context, msg string, kv ...any) {
	logx.Stat(render(msg, l.merge(kv...)))
}

func (l *Logger) merge(kv ...any) []logx.LogField {
	fields := make([]logx.LogField, 0, len(l.fields)+len(kv)/2)
	fields = append(fields, l.fields...)
	return append(fields, Fields(kv...)...)
}

func (l *Logger) warnFields(kv ...any) []logx.LogField {
	return append(l.merge(kv...), logx.Field(LevelKey, WarnLevel))
}

// render 把字段拼进消息，用于 logx 不支持字段的 Severe/Stat。
func render(msg string, fields []logx.LogField) string {
	if len(fields) == 0 {
		return msg
	}
	var b strings.Builder
	b.WriteString(msg)
	for _, f := range fields {
		fmt.Fprintf(&b, " %s=%v", f.Key, f.Value)
	}
	return b.String()
}

// Debug 使用 Default 记录调试日志。
func Debug(ctx context.Context, msg string, kv ...any) {
	Default.Debug(ctx, msg, kv...)
}

// Info 使用 Default 记录普通日志。
func Info(ctx context.Context, msg string, kv ...any) {
	Default.Info(ctx, msg, kv...)
}

// Warn 使用 Default 记录预期内错误。
func Warn(ctx context.Context, msg string, kv ...any) {
	Default.Warn(ctx, msg, kv...)
}

// Error 使用 Default 记录需要及时处理的错误。
func Error(ctx context.Context, msg string, kv ...any) {
	Default.Error(ctx, msg, kv...)
}

// Severe 使用 Default 记录致命日志。
func Severe(ctx context.Context, msg string, kv ...any) {
	Default.Severe(ctx, msg, kv...)
}

// Slow 使用 Default 记录慢日志。
func Slow(ctx context.Context, msg string, kv ...any) {
	Default.Slow(ctx, msg, kv...)
}

// Stat 使用 Default 记录统计日志。
func Stat(ctx context.Context, msg string, kv ...any) {
	Default.Stat(ctx, msg, kv...)
}
