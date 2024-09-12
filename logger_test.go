package pluginlogger_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	// レポジトリ名にハイフンが付いているため、named importを使用する
	pluginlogger "github.com/FlowingSPDG/traefik-plugin-logger"
)

func TestDemo(t *testing.T) {
	cfg := pluginlogger.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := pluginlogger.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	// TODO: assertion...
}
