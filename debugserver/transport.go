package debugserver

import (
	"github.com/eirsyl/flexit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
)

func New(logger log.Logger, healthChecks ...HealthCheck) http.Handler {
	m := http.NewServeMux()
	healthChecker := NewChecker(logger, healthChecks...)

	m.Handle("/healthz", healthChecker.Handler())
	m.Handle("/metrics", promhttp.Handler())

	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return m
}
