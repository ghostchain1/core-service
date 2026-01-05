package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ghostchain1/core-service/internal/config"
	"github.com/ghostchain1/core-service/pkg/logging"
)

func TestGuardAllow(t *testing.T) {
	cfg := &config.Config{Port: "0", LogLevel: "debug"}
	logger := logging.NewLogger(cfg.LogLevel)
	router := NewRouter(cfg, logger)

	body := `{"chain_id":1,"block_number":123,"transactions":10}`
	req := httptest.NewRequest(http.MethodPost, "/guard/op-node", strings.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"allow"`) {
		t.Fatalf("expected allow decision, got: %s", rec.Body.String())
	}
}

func TestGuardDelayLargePayload(t *testing.T) {
	cfg := &config.Config{Port: "0", LogLevel: "debug"}
	logger := logging.NewLogger(cfg.LogLevel)
	router := NewRouter(cfg, logger)

	body := `{"chain_id":1,"block_number":123,"transactions":6000,"calldata_bytes":600000}`
	req := httptest.NewRequest(http.MethodPost, "/guard/proposer", strings.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"delay"`) {
		t.Fatalf("expected delay decision, got: %s", rec.Body.String())
	}
}
