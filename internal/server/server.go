package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ghostchain1/core-service/internal/config"
	"github.com/ghostchain1/core-service/internal/metrics"
	"github.com/ghostchain1/core-service/pkg/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type GuardRequest struct {
	Component      string                 `json:"component"`
	ChainID        uint64                 `json:"chain_id,omitempty"`
	BlockNumber    uint64                 `json:"block_number,omitempty"`
	L1OriginHash   string                 `json:"l1_origin_hash,omitempty"`
	SafeHead       string                 `json:"safe_head,omitempty"`
	FinalizedHead  string                 `json:"finalized_head,omitempty"`
	Transactions   int                    `json:"transactions,omitempty"`
	CalldataBytes  int                    `json:"calldata_bytes,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Identity       string                 `json:"identity,omitempty"`
	RequestID      string                 `json:"request_id,omitempty"`
	TimestampUnix  int64                  `json:"timestamp_unix,omitempty"`
	Source         string                 `json:"source,omitempty"`
	SubmissionHash string                 `json:"submission_hash,omitempty"`
}

type GuardDecision struct {
	Action  string `json:"action"` // allow | delay | block
	Reason  string `json:"reason,omitempty"`
	DelayMS int    `json:"delay_ms,omitempty"`
}

func New(cfg *config.Config, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: NewRouter(cfg, logger),
	}
}

func NewRouter(cfg *config.Config, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(loggingMiddleware(logger))
	r.Use(metrics.Middleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"service": version.ServiceName,
			"version": version.Version,
			"commit":  version.Commit,
			"built":   version.BuildTime,
		})
	})

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/guard", func(r chi.Router) {
		r.Post("/op-node", guardHandler(logger, "op-node"))
		r.Post("/proposer", guardHandler(logger, "proposer"))
	})

	return r
}

func guardHandler(logger *slog.Logger, component string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var req GuardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, logger, http.StatusBadRequest, "invalid_json", err)
			return
		}

		req.Component = component
		if req.TimestampUnix == 0 {
			req.TimestampUnix = time.Now().Unix()
		}

		decision := evaluateGuard(req)
		writeJSON(w, http.StatusOK, decision)
	}
}

func evaluateGuard(req GuardRequest) GuardDecision {
	// Simple scaffold: delay if block is unusually large, otherwise allow.
	if req.CalldataBytes > 500_000 || req.Transactions > 5_000 {
		return GuardDecision{
			Action:  "delay",
			Reason:  "burst_guard_delay",
			DelayMS: 5_000,
		}
	}

	return GuardDecision{
		Action: "allow",
		Reason: "default_allow",
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, logger *slog.Logger, status int, reason string, err error) {
	logger.Warn("request failed", "status", status, "reason", reason, "error", err)
	writeJSON(w, status, map[string]string{
		"error":  reason,
		"detail": err.Error(),
	})
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			route := routePattern(r)
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"route", route,
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"remote", r.RemoteAddr,
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}

func routePattern(r *http.Request) string {
	if ctx := chi.RouteContext(r.Context()); ctx != nil {
		return ctx.RoutePattern()
	}
	return "unknown"
}
