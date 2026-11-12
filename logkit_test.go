package logkit

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
)

// captureWriter 捕获 logx 输出，用于验证转发路径。
type captureWriter struct {
	fields []logx.LogField
	text   string
}

func (w *captureWriter) Alert(any)                            {}
func (w *captureWriter) Close() error                         { return nil }
func (w *captureWriter) Debug(_ any, fields ...logx.LogField) { w.fields = fields }
func (w *captureWriter) Error(_ any, fields ...logx.LogField) { w.fields = fields }
func (w *captureWriter) Info(_ any, fields ...logx.LogField)  { w.fields = fields }
func (w *captureWriter) Severe(v any)                         { w.text = fmt.Sprint(v) }
func (w *captureWriter) Slow(_ any, fields ...logx.LogField)  { w.fields = fields }
func (w *captureWriter) Stack(any)                            {}
func (w *captureWriter) Stat(v any, _ ...logx.LogField)       { w.text = fmt.Sprint(v) }

func capture(t *testing.T) *captureWriter {
	t.Helper()
	w := &captureWriter{}
	logx.SetWriter(w)
	t.Cleanup(func() { logx.SetWriter(logx.NewWriter(os.Stdout)) })
	return w
}

func fieldValueOf(t *testing.T, fields []logx.LogField, key string) any {
	t.Helper()
	for _, f := range fields {
		if f.Key == key {
			return f.Value
		}
	}
	t.Fatalf("缺少字段 %s: %#v", key, fields)
	return nil
}

func TestFieldsOddCount(t *testing.T) {
	fields := Fields("a", 1, "b")
	if len(fields) != 2 {
		t.Fatalf("期望补齐为 2 个字段, 实际 %d", len(fields))
	}
	if fields[1].Key != "b" || fields[1].Value != nil {
		t.Fatalf("奇数字段应补 nil, 实际 %#v", fields[1])
	}
}

func TestFieldsErrorValue(t *testing.T) {
	fields := Fields("err", errors.New("boom"))
	if fields[0].Value != "boom" {
		t.Fatalf("error 应转成字符串, 实际 %#v", fields[0].Value)
	}
}

func TestWarnAppendsBizLevel(t *testing.T) {
	fields := New("service", "xianbaorpc").warnFields("openid", "o1")
	if len(fields) != 3 {
		t.Fatalf("期望 3 个字段, 实际 %d", len(fields))
	}
	if fields[0].Key != "service" || fields[1].Key != "openid" {
		t.Fatalf("公共字段应排在业务字段前, 实际 %#v", fields)
	}
}

func TestWithDoesNotMutateBase(t *testing.T) {
	base := New("a", 1)
	child := base.With("b", 2)
	if len(base.fields) != 1 || len(child.fields) != 2 {
		t.Fatalf("With 不应修改原 logger: base=%d child=%d", len(base.fields), len(child.fields))
	}
}

func TestPackageLevelForwardKV(t *testing.T) {
	w := capture(t)
	Info(context.Background(), "m", "openid", "o1")
	if got := fieldValueOf(t, w.fields, "openid"); got != "o1" {
		t.Fatalf("openid 字段应为 o1, 实际 %#v", got)
	}
}

func TestWarnForwardAddsBizLevel(t *testing.T) {
	w := capture(t)
	Warn(context.Background(), "m", "openid", "o1")
	if got := fieldValueOf(t, w.fields, LevelKey); got != WarnLevel {
		t.Fatalf("%s 应为 %s, 实际 %#v", LevelKey, WarnLevel, got)
	}
	if got := fieldValueOf(t, w.fields, "openid"); got != "o1" {
		t.Fatalf("openid 字段应为 o1, 实际 %#v", got)
	}
}

func TestErrorForward(t *testing.T) {
	w := capture(t)
	Error(context.Background(), "m", "err", errors.New("boom"))
	if got := fieldValueOf(t, w.fields, "err"); got != "boom" {
		t.Fatalf("err 应为 boom, 实际 %#v", got)
	}
}

func TestDebugForward(t *testing.T) {
	w := capture(t)
	logx.SetLevel(logx.DebugLevel)
	t.Cleanup(func() { logx.SetLevel(logx.InfoLevel) })

	Debug(context.Background(), "m", "k", "v")
	if got := fieldValueOf(t, w.fields, "k"); got != "v" {
		t.Fatalf("k 应为 v, 实际 %#v", got)
	}
}

func TestSlowForward(t *testing.T) {
	w := capture(t)
	Slow(context.Background(), "m", "cost", 12)
	if got := fieldValueOf(t, w.fields, "cost"); got != 12 {
		t.Fatalf("cost 应为 12, 实际 %#v", got)
	}
}

func TestSevereAndStatRenderFields(t *testing.T) {
	w := capture(t)

	Severe(context.Background(), "fatal", "a", 1)
	if !strings.Contains(w.text, "a=1") {
		t.Fatalf("Severe 应拼接字段, 实际 %q", w.text)
	}

	w.text = ""
	Stat(context.Background(), "count", "n", 2)
	if !strings.Contains(w.text, "n=2") {
		t.Fatalf("Stat 应拼接字段, 实际 %q", w.text)
	}
}
