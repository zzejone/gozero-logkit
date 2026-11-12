# gozero-logkit

go-zero 项目用的日志包，覆盖 `logx` 的所有级别，并新增 `Warn`。

| 方法 | 级别 | 用途 | 检索/告警方式 |
|---|---|---|---|
| `Debug` | `debug` | 调试，Level=info 时默认不输出 | - |
| `Info` | `info` | 普通日志 | - |
| `Warn` | `info` + `biz_level=warn` | 预期内错误，不需要处理 | 按 `biz_level=warn` 检索 |
| `Error` | `error` | 需要及时处理的错误 | 按 `level=error` 告警 |
| `Severe` | `fatal` | 致命日志 | 按 `level=fatal` 告警 |
| `Slow` | `slow` | 慢日志 | 按 `level=slow` 检索 |
| `Stat` | `stat` | 统计日志 | 按 `level=stat` 检索 |

`Severe`/`Stat` 受 logx 限制不支持结构化字段，字段会拼进消息文本。

> `Warn` 走 info 通道（logx 没有 warn 级别）。`Log.Level` 设为 `error` 时 Warn 不会输出；需要保留 Warn 时 Level 保持 `info`。

## 使用

```go
import "github.com/zzejone/gozero-logkit"

logkit.Info(ctx, "SearchList 入参", "cate_id", 1, "keyword", "面膜")
logkit.Warn(ctx, "攻略已过期", "openid", "o1", "url", "https://...")
logkit.Error(ctx, "转链失败", "openid", "o1", "err", err)
logkit.Debug(ctx, "缓存命中", "key", key)
```

需要公共字段（注入到组件时）：

```go
log := logkit.New("component", "maishou")
log.Info(ctx, "请求完成", "url", url)
log = log.With("openid", openid)
```

- `kv` 是键值对，个数为奇数时最后补 `<nil>`，不会 panic
- `error` 值自动转成字符串
- traceId/spanId 由 `logx.WithContext` 自动带上（`Severe`/`Stat` 除外）

## 版本

```bash
git tag v0.1.0 && git push origin v0.1.0
go get github.com/zzejone/gozero-logkit@v0.1.0
```
