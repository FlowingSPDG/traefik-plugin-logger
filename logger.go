// Package pluginlogger is a demo logger plugin.
package pluginlogger

// packageのimport
// ここで標準パッケージ以外を用いる場合、`go mod vendor` が必須となる
// 基本的には標準パッケージのみで開発することをオススメする
import (
	"context"
	"io"
	"net/http"
	"os"
)

// Config プラグインの設定が入った構造体
// yaegiから設定情報が注入されるので、全フィールドをPublicにする
type Config struct {
	Prefix string `json:"prefix,omitempty"`
}

// CreateConfig 設定を作成する関数
// yaegiから呼び出されるため、必須
func CreateConfig() *Config {
	return &Config{
		Prefix: "[YOUR_PLUGIN]",
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
	out io.Writer
}

// New プラグイン本体を作成する関数
// CreateConfig() と同じく、Yaegiから直接呼び出されるので必須
// また、関数シグネチャは`New(context.Context, http.Handler, *Config, string) (http.Handler, error)` である必要がある
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Logger{
		next: next,
		name: name,
		cfg:  config,
		out:  os.Stdout,
	}, nil
}

// ServeHTTP はミドルウェアの実処理
func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO: 実装
	l.next.ServeHTTP(rw, req)
}
