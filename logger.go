// Package pluginlogger is a demo logger plugin.
package pluginlogger

// packageのimport
// ここで標準パッケージ以外を用いる場合、`go mod vendor` が必須となる
// 基本的には標準パッケージのみで開発することをオススメする
import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

// Config プラグインの設定が入った構造体
// yaegiから設定情報が注入されるので、全フィールドをPublicにする
type Config struct {
	Prefix   string `json:"prefix,omitempty"`
	LogLevel string `json:"loglevel,omitempty"`
}

// logLevel stringのLogLevelをslog.Level に変換する
func (c *Config) logLevel() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// CreateConfig 設定を作成する関数
// yaegiから呼び出されるため、必須
func CreateConfig() *Config {
	return &Config{
		Prefix:   "[YOUR_PLUGIN]",
		LogLevel: "info",
	}
}

// Logger プラグインの本体。必要な依存性を内部に保持する
// http.Handler を実装する
type Logger struct {
	// next, name はほぼ必須フィールド
	next http.Handler
	name string

	// 設定はそのまま保持する
	cfg *Config

	// ログの出力先
	logger *slog.Logger
}

// New プラグイン本体を作成する関数
// CreateConfig() と同じく、Yaegiから直接呼び出されるので必須
// また、関数シグネチャは`New(context.Context, http.Handler, *Config, string) (http.Handler, error)` である必要がある
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	out := os.Stdout
	return &Logger{
		next:   next,
		name:   name,
		cfg:    config,
		logger: slog.New(slog.NewTextHandler(out, nil)),
	}, nil
}

// convertRequest *http.Request をslog.Attr に変換するユーティリティ関数
func convertRequest(req *http.Request) []slog.Attr {
	return []slog.Attr{
		slog.String("proto", req.Proto),
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
	}
}

// ServeHTTP はミドルウェアの実処理
func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// contextを取得
	ctx := req.Context()

	// fieldに変換
	field := convertRequest(req)

	// ログを出力
	l.logger.LogAttrs(ctx, l.cfg.logLevel(), l.cfg.Prefix, field...)

	// next...
	l.next.ServeHTTP(rw, req)

	// 後処理
}
